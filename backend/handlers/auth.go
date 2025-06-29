package handlers

import (
	"database/sql"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"arvfinder-backend/database"
	"arvfinder-backend/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	authService   *services.AuthService
	rateLimiter   *services.RateLimiter
	sms2FAService *services.SMS2FAService
	db            *sql.DB
}

// LoginResponse represents the response for successful login
type LoginResponse struct {
	Success      bool                    `json:"success"`
	Message      string                  `json:"message"`
	User         *services.User          `json:"user,omitempty"`
	Tokens       *services.TokenPair     `json:"tokens,omitempty"`
	Requires2FA  bool                    `json:"requires_2fa"`
	TempToken    string                  `json:"temp_token,omitempty"` // For 2FA flow
}

// RegisterResponse represents the response for registration
type RegisterResponse struct {
	Success           bool   `json:"success"`
	Message           string `json:"message"`
	UserID            string `json:"user_id,omitempty"`
	RequiresVerification bool `json:"requires_verification"`
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler() *AuthHandler {
	// Get database connection
	db := database.GetDB()
	
	// Get JWT secret from environment
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-super-secret-jwt-key-change-in-production" // Default for development
	}

	// Initialize services
	authService := services.NewAuthService(db, jwtSecret)
	rateLimiter := services.NewRateLimiter(db)
	
	// Initialize SMS 2FA service
	twilioSID := os.Getenv("TWILIO_ACCOUNT_SID")
	twilioToken := os.Getenv("TWILIO_AUTH_TOKEN")
	twilioPhone := os.Getenv("TWILIO_PHONE_NUMBER")
	sms2FAService := services.NewSMS2FAService(db, authService, twilioSID, twilioToken, twilioPhone)

	return &AuthHandler{
		authService:   authService,
		rateLimiter:   rateLimiter,
		sms2FAService: sms2FAService,
		db:            db,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	clientIP := h.getClientIP(c)
	userAgent := c.GetHeader("User-Agent")

	// Check rate limiting
	allowed, blockTime, err := h.rateLimiter.IsAllowed(clientIP, "register")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal server error",
		})
		return
	}

	if !allowed {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"success":       false,
			"message":       "Too many registration attempts. Please try again later.",
			"retry_after":   int(blockTime.Seconds()),
		})
		return
	}

	// Record the attempt
	h.rateLimiter.RecordAttempt(clientIP, "register")

	var req services.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
			"errors":  err.Error(),
		})
		return
	}

	// Validate password strength
	if !h.isPasswordStrong(req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Password must be at least 8 characters with uppercase, lowercase, number, and special character",
		})
		return
	}

	// Check if user already exists
	var existingUserID string
	err = h.db.QueryRow("SELECT id FROM users WHERE email = $1", req.Email).Scan(&existingUserID)
	if err != sql.ErrNoRows {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": "User with this email already exists",
		})
		return
	}

	// Create or get tenant
	tenantID, err := h.createOrGetTenant(req.TenantName, req.Email)
	if err != nil {
		h.authService.LogSecurityEvent("", "registration_failed", "Failed to create tenant", clientIP, userAgent, map[string]interface{}{
			"email": req.Email,
			"error": err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create account",
		})
		return
	}

	// Generate salt and hash password
	salt, err := h.authService.GenerateSecureSalt()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create account",
		})
		return
	}

	passwordHash := h.authService.HashPassword(req.Password, salt)
	saltString := string(salt)

	// Generate email verification token
	emailVerificationToken := uuid.New().String()
	emailVerificationExpires := time.Now().Add(24 * time.Hour)

	// Create user
	userID := uuid.New().String()
	_, err = h.db.Exec(`
		INSERT INTO users (
			id, tenant_id, email, password_hash, password_salt, first_name, last_name, 
			phone_number, email_verification_token, email_verification_expires_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`, userID, tenantID, req.Email, passwordHash, saltString, req.FirstName, req.LastName, 
		req.PhoneNumber, emailVerificationToken, emailVerificationExpires)

	if err != nil {
		h.authService.LogSecurityEvent("", "registration_failed", "Database error during user creation", clientIP, userAgent, map[string]interface{}{
			"email": req.Email,
			"error": err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create account",
		})
		return
	}

	// Reset rate limiter on successful registration
	h.rateLimiter.ResetAttempts(clientIP, "register")

	// Log successful registration
	h.authService.LogSecurityEvent(userID, "user_registered", "User successfully registered", clientIP, userAgent, map[string]interface{}{
		"email": req.Email,
	})

	// In production, you would send an email verification here
	// For now, we'll just return success

	c.JSON(http.StatusCreated, RegisterResponse{
		Success:              true,
		Message:              "Account created successfully. Please verify your email address.",
		UserID:               userID,
		RequiresVerification: true,
	})
}

// Login handles user authentication
func (h *AuthHandler) Login(c *gin.Context) {
	clientIP := h.getClientIP(c)
	userAgent := c.GetHeader("User-Agent")

	// Check rate limiting
	allowed, blockTime, err := h.rateLimiter.IsAllowed(clientIP, "login")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal server error",
		})
		return
	}

	if !allowed {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"success":     false,
			"message":     "Too many login attempts. Please try again later.",
			"retry_after": int(blockTime.Seconds()),
		})
		return
	}

	// Record the attempt
	h.rateLimiter.RecordAttempt(clientIP, "login")

	var req services.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
		})
		return
	}

	req.IPAddress = clientIP
	req.DeviceInfo = userAgent

	// Get user from database
	var user services.User
	var passwordHash, passwordSalt string
	err = h.db.QueryRow(`
		SELECT id, tenant_id, email, password_hash, password_salt, first_name, last_name, 
		       phone_number, phone_verified, role, is_active, two_factor_enabled, 
		       last_login_at, failed_login_attempts, locked_until, created_at, updated_at, email_verified
		FROM users WHERE email = $1
	`, req.Email).Scan(
		&user.ID, &user.TenantID, &user.Email, &passwordHash, &passwordSalt,
		&user.FirstName, &user.LastName, &user.PhoneNumber, &user.PhoneVerified,
		&user.Role, &user.IsActive, &user.TwoFactorEnabled, &user.LastLoginAt,
		&user.FailedLoginAttempts, &user.LockedUntil, &user.CreatedAt, &user.UpdatedAt, &user.EmailVerified,
	)

	if err == sql.ErrNoRows {
		h.authService.LogSecurityEvent("", "login_failed", "User not found", clientIP, userAgent, map[string]interface{}{
			"email": req.Email,
		})
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid email or password",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal server error",
		})
		return
	}

	// Check if account is active
	if !user.IsActive {
		h.authService.LogSecurityEvent(user.ID, "login_failed", "Account inactive", clientIP, userAgent, nil)
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Account is inactive. Please contact support.",
		})
		return
	}

	// Check if account is locked
	if user.LockedUntil != nil && user.LockedUntil.After(time.Now()) {
		remaining := time.Until(*user.LockedUntil)
		h.authService.LogSecurityEvent(user.ID, "login_failed", "Account locked", clientIP, userAgent, map[string]interface{}{
			"locked_until": user.LockedUntil,
		})
		c.JSON(http.StatusUnauthorized, gin.H{
			"success":     false,
			"message":     "Account is temporarily locked due to too many failed attempts",
			"retry_after": int(remaining.Seconds()),
		})
		return
	}

	// Verify password
	if !h.authService.VerifyPassword(req.Password, passwordHash) {
		// Increment failed attempts
		h.authService.IncrementFailedAttempts(user.ID)
		h.authService.LogSecurityEvent(user.ID, "login_failed", "Invalid password", clientIP, userAgent, nil)
		
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid email or password",
		})
		return
	}

	// Check if email is verified
	if !user.EmailVerified {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Please verify your email address before logging in",
		})
		return
	}

	// Check if 2FA is enabled
	if user.TwoFactorEnabled && user.PhoneVerified {
		// Send 2FA code
		smsRequest := &services.SMSVerificationRequest{
			PhoneNumber: user.PhoneNumber,
			Purpose:     "login",
			UserID:      user.ID,
		}

		_, err = h.sms2FAService.SendVerificationCode(smsRequest)
		if err != nil {
			h.authService.LogSecurityEvent(user.ID, "2fa_send_failed", "Failed to send 2FA code", clientIP, userAgent, map[string]interface{}{
				"error": err.Error(),
			})
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to send verification code",
			})
			return
		}

		// Generate temporary token for 2FA flow
		tempToken := uuid.New().String()
		
		c.JSON(http.StatusOK, LoginResponse{
			Success:     true,
			Message:     "Verification code sent to your phone",
			Requires2FA: true,
			TempToken:   tempToken,
		})
		return
	}

	// Generate token pair
	tokens, err := h.authService.GenerateTokenPair(&user, req.DeviceInfo, req.IPAddress)
	if err != nil {
		h.authService.LogSecurityEvent(user.ID, "login_failed", "Token generation failed", clientIP, userAgent, map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to generate authentication tokens",
		})
		return
	}

	// Reset failed attempts and update last login
	h.authService.ResetFailedAttempts(user.ID)
	h.rateLimiter.ResetAttempts(clientIP, "login")

	// Log successful login
	h.authService.LogSecurityEvent(user.ID, "login_success", "User successfully logged in", clientIP, userAgent, nil)

	// Remove sensitive data from user object
	user.FailedLoginAttempts = 0
	user.LockedUntil = nil

	c.JSON(http.StatusOK, LoginResponse{
		Success:     true,
		Message:     "Login successful",
		User:        &user,
		Tokens:      tokens,
		Requires2FA: false,
	})
}

// Verify2FA handles 2FA code verification during login
func (h *AuthHandler) Verify2FA(c *gin.Context) {
	clientIP := h.getClientIP(c)
	userAgent := c.GetHeader("User-Agent")

	var req services.VerifyCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
		})
		return
	}

	// Verify the 2FA code
	response, err := h.sms2FAService.VerifyCode(&req)
	if err != nil {
		h.authService.LogSecurityEvent(req.UserID, "2fa_verification_failed", "2FA verification error", clientIP, userAgent, map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Verification failed",
		})
		return
	}

	if !response.Verified {
		h.authService.LogSecurityEvent(req.UserID, "2fa_verification_failed", "Invalid 2FA code", clientIP, userAgent, nil)
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": response.Message,
		})
		return
	}

	// Get user for token generation
	var user services.User
	err = h.db.QueryRow(`
		SELECT id, tenant_id, email, first_name, last_name, phone_number, 
		       phone_verified, role, is_active, two_factor_enabled, created_at, updated_at
		FROM users WHERE id = $1 AND is_active = TRUE
	`, req.UserID).Scan(
		&user.ID, &user.TenantID, &user.Email, &user.FirstName, &user.LastName,
		&user.PhoneNumber, &user.PhoneVerified, &user.Role, &user.IsActive,
		&user.TwoFactorEnabled, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User not found or inactive",
		})
		return
	}

	// Generate token pair
	tokens, err := h.authService.GenerateTokenPair(&user, userAgent, clientIP)
	if err != nil {
		h.authService.LogSecurityEvent(user.ID, "2fa_login_failed", "Token generation failed after 2FA", clientIP, userAgent, map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to complete login",
		})
		return
	}

	// Reset failed attempts and update last login
	h.authService.ResetFailedAttempts(user.ID)
	h.rateLimiter.ResetAttempts(clientIP, "login")

	// Log successful 2FA login
	h.authService.LogSecurityEvent(user.ID, "2fa_login_success", "User successfully logged in with 2FA", clientIP, userAgent, nil)

	c.JSON(http.StatusOK, LoginResponse{
		Success:     true,
		Message:     "Login successful",
		User:        &user,
		Tokens:      tokens,
		Requires2FA: false,
	})
}

// Helper functions

func (h *AuthHandler) getClientIP(c *gin.Context) string {
	// Check various headers for the real IP
	ip := c.GetHeader("X-Forwarded-For")
	if ip == "" {
		ip = c.GetHeader("X-Real-IP")
	}
	if ip == "" {
		ip = c.ClientIP()
	}
	
	// Extract first IP if comma-separated
	if strings.Contains(ip, ",") {
		ip = strings.TrimSpace(strings.Split(ip, ",")[0])
	}
	
	// Validate IP address
	if net.ParseIP(ip) == nil {
		ip = "127.0.0.1" // Fallback
	}
	
	return ip
}

func (h *AuthHandler) isPasswordStrong(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasDigit = true
		case strings.ContainsRune("!@#$%^&*()_+-=[]{}|;:,.<>?", char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasDigit && hasSpecial
}

func (h *AuthHandler) createOrGetTenant(tenantName, userEmail string) (string, error) {
	// If no tenant name provided, create a personal tenant
	if tenantName == "" {
		tenantName = "Personal Account"
	}

	tenantID := uuid.New().String()
	_, err := h.db.Exec(`
		INSERT INTO tenants (id, name, subscription_tier) 
		VALUES ($1, $2, 'starter')
	`, tenantID, tenantName)

	if err != nil {
		return "", err
	}

	return tenantID, nil
}

