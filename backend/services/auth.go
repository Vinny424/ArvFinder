package services

import (
	"crypto/rand"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

// AuthService handles user authentication with extreme security measures
type AuthService struct {
	db             *sql.DB
	jwtSecret      []byte
	argon2Params   *Argon2Params
	tokenDuration  time.Duration
	refreshDuration time.Duration
}

// Argon2Params defines parameters for Argon2 password hashing
type Argon2Params struct {
	Memory      uint32 // Memory in KB
	Iterations  uint32 // Number of iterations
	Parallelism uint8  // Number of threads
	SaltLength  uint32 // Salt length in bytes
	KeyLength   uint32 // Key length in bytes
}

// User represents a user in the system
type User struct {
	ID                    string     `json:"id"`
	TenantID              string     `json:"tenant_id"`
	Email                 string     `json:"email"`
	EmailVerified         bool       `json:"email_verified"`
	FirstName             string     `json:"first_name,omitempty"`
	LastName              string     `json:"last_name,omitempty"`
	PhoneNumber           string     `json:"phone_number,omitempty"`
	PhoneVerified         bool       `json:"phone_verified"`
	Role                  string     `json:"role"`
	IsActive              bool       `json:"is_active"`
	TwoFactorEnabled      bool       `json:"two_factor_enabled"`
	LastLoginAt           *time.Time `json:"last_login_at,omitempty"`
	FailedLoginAttempts   int        `json:"failed_login_attempts"`
	LockedUntil           *time.Time `json:"locked_until,omitempty"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

// TokenPair represents access and refresh tokens
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	TokenType    string    `json:"token_type"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	RememberMe  bool   `json:"remember_me"`
	DeviceInfo  string `json:"device_info,omitempty"`
	IPAddress   string `json:"ip_address,omitempty"`
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	FirstName   string `json:"first_name" binding:"required,min=1,max=100"`
	LastName    string `json:"last_name" binding:"required,min=1,max=100"`
	PhoneNumber string `json:"phone_number,omitempty"`
	TenantName  string `json:"tenant_name,omitempty"`
}

// JWTClaims represents JWT token claims
type JWTClaims struct {
	UserID       string `json:"user_id"`
	TenantID     string `json:"tenant_id"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	SessionID    string `json:"session_id"`
	DeviceFingerprint string `json:"device_fingerprint,omitempty"`
	jwt.RegisteredClaims
}

// NewAuthService creates a new authentication service
func NewAuthService(db *sql.DB, jwtSecret string) *AuthService {
	// Production-grade Argon2 parameters
	argon2Params := &Argon2Params{
		Memory:      128 * 1024, // 128 MB
		Iterations:  4,          // 4 iterations
		Parallelism: 4,          // 4 threads
		SaltLength:  32,         // 32 bytes salt
		KeyLength:   64,         // 64 bytes key
	}

	return &AuthService{
		db:             db,
		jwtSecret:      []byte(jwtSecret),
		argon2Params:   argon2Params,
		tokenDuration:  15 * time.Minute,  // Access token: 15 minutes
		refreshDuration: 7 * 24 * time.Hour, // Refresh token: 7 days
	}
}

// GenerateSecureSalt generates a cryptographically secure random salt
func (a *AuthService) GenerateSecureSalt() ([]byte, error) {
	salt := make([]byte, a.argon2Params.SaltLength)
	_, err := rand.Read(salt)
	return salt, err
}

// HashPassword creates a secure Argon2 hash of the password
func (a *AuthService) HashPassword(password string, salt []byte) string {
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		a.argon2Params.Iterations,
		a.argon2Params.Memory,
		a.argon2Params.Parallelism,
		a.argon2Params.KeyLength,
	)
	
	// Format: $argon2id$v=19$m=128,t=4,p=4$salt$hash
	return fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		a.argon2Params.Memory,
		a.argon2Params.Iterations,
		a.argon2Params.Parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)
}

// VerifyPassword verifies a password against its hash using constant-time comparison
func (a *AuthService) VerifyPassword(password, hashedPassword string) bool {
	// Parse the hash format
	parts := strings.Split(hashedPassword, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return false
	}

	// Extract salt and hash
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false
	}
	
	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false
	}

	// Compute hash of provided password
	actualHash := argon2.IDKey(
		[]byte(password),
		salt,
		a.argon2Params.Iterations,
		a.argon2Params.Memory,
		a.argon2Params.Parallelism,
		a.argon2Params.KeyLength,
	)

	// Use constant-time comparison to prevent timing attacks
	return subtle.ConstantTimeCompare(expectedHash, actualHash) == 1
}

// GenerateTokenPair creates a new access/refresh token pair
func (a *AuthService) GenerateTokenPair(user *User, deviceInfo, ipAddress string) (*TokenPair, error) {
	// Generate session ID
	sessionID := uuid.New().String()
	
	// Create device fingerprint (simplified)
	deviceFingerprint := a.createDeviceFingerprint(deviceInfo, ipAddress)
	
	// Create access token claims
	accessClaims := &JWTClaims{
		UserID:            user.ID,
		TenantID:          user.TenantID,
		Email:             user.Email,
		Role:              user.Role,
		SessionID:         sessionID,
		DeviceFingerprint: deviceFingerprint,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Subject:   user.ID,
			Issuer:    "arvfinder",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.tokenDuration)),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Generate access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(a.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshTokenBytes := make([]byte, 32)
	_, err = rand.Read(refreshTokenBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	refreshToken := base64.URLEncoding.EncodeToString(refreshTokenBytes)
	
	// Hash refresh token for storage
	refreshSalt, err := a.GenerateSecureSalt()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token salt: %w", err)
	}
	refreshTokenHash := a.HashPassword(refreshToken, refreshSalt)

	// Store session in database
	expiresAt := time.Now().Add(a.refreshDuration)
	_, err = a.db.Exec(`
		INSERT INTO user_sessions (
			user_id, refresh_token, refresh_token_hash, access_token_jti,
			device_fingerprint, user_agent, ip_address, expires_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		user.ID, refreshToken, refreshTokenHash, accessClaims.ID,
		deviceFingerprint, deviceInfo, ipAddress, expiresAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to store session: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		TokenType:    "Bearer",
	}, nil
}

// ValidateToken validates and parses a JWT token
func (a *AuthService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		// Check if session is still valid
		var sessionExists bool
		err = a.db.QueryRow(`
			SELECT EXISTS(
				SELECT 1 FROM user_sessions 
				WHERE access_token_jti = $1 AND expires_at > NOW() AND revoked = FALSE
			)`, claims.ID).Scan(&sessionExists)
		
		if err != nil || !sessionExists {
			return nil, fmt.Errorf("session invalid or expired")
		}
		
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// createDeviceFingerprint creates a simple device fingerprint
func (a *AuthService) createDeviceFingerprint(deviceInfo, ipAddress string) string {
	// Simplified fingerprinting - in production, you'd use more sophisticated methods
	fingerprint := fmt.Sprintf("%s:%s", deviceInfo, ipAddress)
	return base64.StdEncoding.EncodeToString([]byte(fingerprint))
}

// ValidateIPAddress validates an IP address format
func (a *AuthService) ValidateIPAddress(ip string) bool {
	return net.ParseIP(ip) != nil
}

// IsAccountLocked checks if a user account is currently locked
func (a *AuthService) IsAccountLocked(userID string) (bool, time.Duration, error) {
	var lockedUntil *time.Time
	err := a.db.QueryRow(`
		SELECT locked_until FROM users WHERE id = $1
	`, userID).Scan(&lockedUntil)
	
	if err != nil {
		return false, 0, err
	}
	
	if lockedUntil != nil && lockedUntil.After(time.Now()) {
		remaining := time.Until(*lockedUntil)
		return true, remaining, nil
	}
	
	return false, 0, nil
}

// IncrementFailedAttempts increments the failed login attempts counter
func (a *AuthService) IncrementFailedAttempts(userID string) error {
	_, err := a.db.Exec(`
		UPDATE users 
		SET failed_login_attempts = failed_login_attempts + 1,
		    updated_at = NOW()
		WHERE id = $1
	`, userID)
	return err
}

// ResetFailedAttempts resets the failed login attempts counter
func (a *AuthService) ResetFailedAttempts(userID string) error {
	_, err := a.db.Exec(`
		UPDATE users 
		SET failed_login_attempts = 0,
		    locked_until = NULL,
		    last_login_at = NOW(),
		    updated_at = NOW()
		WHERE id = $1
	`, userID)
	return err
}

// LogSecurityEvent logs a security event to the audit log
func (a *AuthService) LogSecurityEvent(userID, eventType, description, ipAddress, userAgent string, additionalData map[string]interface{}) error {
	var jsonData interface{}
	if additionalData != nil {
		// Convert map to JSON (simplified - in production use proper JSON handling)
		jsonData = additionalData
	}
	
	_, err := a.db.Exec(`
		INSERT INTO security_audit_log (
			user_id, event_type, event_description, ip_address, user_agent, additional_data
		) VALUES ($1, $2, $3, $4, $5, $6)`,
		userID, eventType, description, ipAddress, userAgent, jsonData,
	)
	return err
}

// RevokeSession revokes a user session
func (a *AuthService) RevokeSession(refreshToken string) error {
	_, err := a.db.Exec(`
		UPDATE user_sessions 
		SET revoked = TRUE 
		WHERE refresh_token = $1
	`, refreshToken)
	return err
}

// RevokeAllUserSessions revokes all sessions for a user
func (a *AuthService) RevokeAllUserSessions(userID string) error {
	_, err := a.db.Exec(`
		UPDATE user_sessions 
		SET revoked = TRUE 
		WHERE user_id = $1
	`, userID)
	return err
}

// CleanupExpiredSessions removes expired sessions from the database
func (a *AuthService) CleanupExpiredSessions() error {
	_, err := a.db.Exec(`DELETE FROM user_sessions WHERE expires_at < NOW()`)
	return err
}