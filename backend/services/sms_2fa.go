package services

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// SMS2FAService handles SMS-based two-factor authentication
type SMS2FAService struct {
	db           *sql.DB
	authService  *AuthService
	twilioSID    string
	twilioToken  string
	twilioPhone  string
	testMode     bool // For testing without actual SMS
}

// SMSVerificationRequest represents an SMS verification request
type SMSVerificationRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	Purpose     string `json:"purpose" binding:"required"` // 'login', 'register', 'password_reset'
	UserID      string `json:"user_id,omitempty"`
}

// SMSVerificationResponse represents the response to SMS verification request
type SMSVerificationResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	CodeSent  bool   `json:"code_sent"`
	ExpiresAt int64  `json:"expires_at"`
}

// VerifyCodeRequest represents a code verification request
type VerifyCodeRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	Code        string `json:"code" binding:"required,len=6"`
	Purpose     string `json:"purpose" binding:"required"`
	UserID      string `json:"user_id,omitempty"`
}

// VerifyCodeResponse represents the response to code verification
type VerifyCodeResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	Verified bool   `json:"verified"`
}

// NewSMS2FAService creates a new SMS 2FA service
// For Twilio integration, you'll need to provide:
// - Twilio Account SID
// - Twilio Auth Token  
// - Twilio Phone Number
func NewSMS2FAService(db *sql.DB, authService *AuthService, twilioSID, twilioToken, twilioPhone string) *SMS2FAService {
	testMode := twilioSID == "" || twilioToken == "" || twilioPhone == ""
	
	return &SMS2FAService{
		db:          db,
		authService: authService,
		twilioSID:   twilioSID,
		twilioToken: twilioToken,
		twilioPhone: twilioPhone,
		testMode:    testMode,
	}
}

// GenerateVerificationCode generates a secure 6-digit verification code
func (s *SMS2FAService) GenerateVerificationCode() (string, error) {
	// Generate a secure random 6-digit code
	max := big.NewInt(999999)
	min := big.NewInt(100000)
	
	n, err := rand.Int(rand.Reader, max.Sub(max, min).Add(max, big.NewInt(1)))
	if err != nil {
		return "", fmt.Errorf("failed to generate random code: %w", err)
	}
	
	code := n.Add(n, min).String()
	return code, nil
}

// SendVerificationCode sends a verification code via SMS
func (s *SMS2FAService) SendVerificationCode(request *SMSVerificationRequest) (*SMSVerificationResponse, error) {
	// Validate phone number format (basic validation)
	if !s.isValidPhoneNumber(request.PhoneNumber) {
		return &SMSVerificationResponse{
			Success: false,
			Message: "Invalid phone number format",
		}, nil
	}

	// Generate verification code
	code, err := s.GenerateVerificationCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate verification code: %w", err)
	}

	// Hash the code for storage
	salt, err := s.authService.GenerateSecureSalt()
	if err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	codeHash := s.authService.HashPassword(code, salt)

	// Set expiration (5 minutes from now)
	expiresAt := time.Now().Add(5 * time.Minute)

	// Clean up any existing unexpired codes for this phone/purpose
	_, err = s.db.Exec(`
		DELETE FROM sms_verification_codes 
		WHERE phone_number = $1 AND purpose = $2 AND expires_at > NOW()
	`, request.PhoneNumber, request.Purpose)
	if err != nil {
		return nil, fmt.Errorf("failed to cleanup existing codes: %w", err)
	}

	// Store the verification code
	_, err = s.db.Exec(`
		INSERT INTO sms_verification_codes (
			user_id, phone_number, code, code_hash, purpose, expires_at
		) VALUES ($1, $2, $3, $4, $5, $6)
	`, request.UserID, request.PhoneNumber, code, codeHash, request.Purpose, expiresAt)
	if err != nil {
		return nil, fmt.Errorf("failed to store verification code: %w", err)
	}

	// Send SMS
	var smsSent bool
	if s.testMode {
		// In test mode, log the code instead of sending SMS
		fmt.Printf("TEST MODE: SMS verification code for %s: %s\n", request.PhoneNumber, code)
		smsSent = true
	} else {
		smsSent, err = s.sendSMSViaTwilio(request.PhoneNumber, code, request.Purpose)
		if err != nil {
			return nil, fmt.Errorf("failed to send SMS: %w", err)
		}
	}

	return &SMSVerificationResponse{
		Success:   true,
		Message:   "Verification code sent successfully",
		CodeSent:  smsSent,
		ExpiresAt: expiresAt.Unix(),
	}, nil
}

// VerifyCode verifies a submitted verification code
func (s *SMS2FAService) VerifyCode(request *VerifyCodeRequest) (*VerifyCodeResponse, error) {
	// Get the stored verification record
	var storedCode, codeHash string
	var attempts, maxAttempts int
	var expiresAt time.Time
	var verified bool

	err := s.db.QueryRow(`
		SELECT code, code_hash, attempts, max_attempts, expires_at, verified
		FROM sms_verification_codes 
		WHERE phone_number = $1 AND purpose = $2 
		ORDER BY created_at DESC 
		LIMIT 1
	`, request.PhoneNumber, request.Purpose).Scan(
		&storedCode, &codeHash, &attempts, &maxAttempts, &expiresAt, &verified,
	)

	if err == sql.ErrNoRows {
		return &VerifyCodeResponse{
			Success:  false,
			Message:  "No verification code found",
			Verified: false,
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get verification code: %w", err)
	}

	// Check if already verified
	if verified {
		return &VerifyCodeResponse{
			Success:  false,
			Message:  "Code has already been used",
			Verified: false,
		}, nil
	}

	// Check if expired
	if time.Now().After(expiresAt) {
		return &VerifyCodeResponse{
			Success:  false,
			Message:  "Verification code has expired",
			Verified: false,
		}, nil
	}

	// Check if too many attempts
	if attempts >= maxAttempts {
		return &VerifyCodeResponse{
			Success:  false,
			Message:  "Too many verification attempts. Please request a new code.",
			Verified: false,
		}, nil
	}

	// Increment attempts
	_, err = s.db.Exec(`
		UPDATE sms_verification_codes 
		SET attempts = attempts + 1 
		WHERE phone_number = $1 AND purpose = $2 AND expires_at = $3
	`, request.PhoneNumber, request.Purpose, expiresAt)
	if err != nil {
		return nil, fmt.Errorf("failed to increment attempts: %w", err)
	}

	// Verify the code using constant-time comparison
	isValid := s.authService.VerifyPassword(request.Code, codeHash)

	if !isValid {
		remainingAttempts := maxAttempts - (attempts + 1)
		message := fmt.Sprintf("Invalid verification code. %d attempts remaining.", remainingAttempts)
		if remainingAttempts <= 0 {
			message = "Invalid verification code. No more attempts allowed. Please request a new code."
		}
		
		return &VerifyCodeResponse{
			Success:  false,
			Message:  message,
			Verified: false,
		}, nil
	}

	// Mark as verified
	_, err = s.db.Exec(`
		UPDATE sms_verification_codes 
		SET verified = TRUE 
		WHERE phone_number = $1 AND purpose = $2 AND expires_at = $3
	`, request.PhoneNumber, request.Purpose, expiresAt)
	if err != nil {
		return nil, fmt.Errorf("failed to mark code as verified: %w", err)
	}

	// If this is for phone verification, update user's phone_verified status
	if request.Purpose == "phone_verification" && request.UserID != "" {
		_, err = s.db.Exec(`
			UPDATE users 
			SET phone_verified = TRUE, phone_number = $1, updated_at = NOW() 
			WHERE id = $2
		`, request.PhoneNumber, request.UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to update user phone verification: %w", err)
		}
	}

	return &VerifyCodeResponse{
		Success:  true,
		Message:  "Verification code verified successfully",
		Verified: true,
	}, nil
}

// sendSMSViaTwilio sends SMS using Twilio API
func (s *SMS2FAService) sendSMSViaTwilio(phoneNumber, code, purpose string) (bool, error) {
	if s.testMode {
		return false, fmt.Errorf("Twilio not configured - running in test mode")
	}

	// Format the message based on purpose
	var message string
	switch purpose {
	case "login":
		message = fmt.Sprintf("Your ArvFinder login verification code is: %s. This code expires in 5 minutes.", code)
	case "register":
		message = fmt.Sprintf("Your ArvFinder registration verification code is: %s. This code expires in 5 minutes.", code)
	case "password_reset":
		message = fmt.Sprintf("Your ArvFinder password reset code is: %s. This code expires in 5 minutes.", code)
	case "phone_verification":
		message = fmt.Sprintf("Your ArvFinder phone verification code is: %s. This code expires in 5 minutes.", code)
	default:
		message = fmt.Sprintf("Your ArvFinder verification code is: %s. This code expires in 5 minutes.", code)
	}

	// Prepare Twilio API request
	apiURL := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", s.twilioSID)
	
	data := url.Values{}
	data.Set("From", s.twilioPhone)
	data.Set("To", phoneNumber)
	data.Set("Body", message)

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(s.twilioSID, s.twilioToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return true, nil
	}

	return false, fmt.Errorf("Twilio API returned status: %d", resp.StatusCode)
}

// isValidPhoneNumber performs basic phone number validation
func (s *SMS2FAService) isValidPhoneNumber(phone string) bool {
	// Remove common formatting characters
	cleaned := strings.ReplaceAll(phone, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, "(", "")
	cleaned = strings.ReplaceAll(cleaned, ")", "")
	cleaned = strings.ReplaceAll(cleaned, ".", "")

	// Should start with + and have 10-15 digits
	if !strings.HasPrefix(cleaned, "+") {
		return false
	}

	digits := cleaned[1:] // Remove the +
	if len(digits) < 10 || len(digits) > 15 {
		return false
	}

	// Check if all characters after + are digits
	for _, char := range digits {
		if char < '0' || char > '9' {
			return false
		}
	}

	return true
}

// CleanupExpiredCodes removes expired verification codes
func (s *SMS2FAService) CleanupExpiredCodes() error {
	_, err := s.db.Exec(`
		DELETE FROM sms_verification_codes 
		WHERE expires_at < NOW() - INTERVAL '1 day'
	`)
	return err
}

// GetVerificationStatus returns the status of verification for a phone number
func (s *SMS2FAService) GetVerificationStatus(phoneNumber, purpose string) (bool, time.Time, int, error) {
	var verified bool
	var expiresAt time.Time
	var attempts int

	err := s.db.QueryRow(`
		SELECT verified, expires_at, attempts
		FROM sms_verification_codes 
		WHERE phone_number = $1 AND purpose = $2 
		ORDER BY created_at DESC 
		LIMIT 1
	`, phoneNumber, purpose).Scan(&verified, &expiresAt, &attempts)

	if err == sql.ErrNoRows {
		return false, time.Time{}, 0, nil
	}
	if err != nil {
		return false, time.Time{}, 0, fmt.Errorf("failed to get verification status: %w", err)
	}

	return verified, expiresAt, attempts, nil
}

// RevokeVerificationCode revokes an unused verification code
func (s *SMS2FAService) RevokeVerificationCode(phoneNumber, purpose string) error {
	_, err := s.db.Exec(`
		DELETE FROM sms_verification_codes 
		WHERE phone_number = $1 AND purpose = $2 AND verified = FALSE
	`, phoneNumber, purpose)
	return err
}