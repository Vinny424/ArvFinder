package services

import (
	"database/sql"
	"fmt"
	"time"
)

// RateLimiter handles rate limiting and brute force protection
type RateLimiter struct {
	db *sql.DB
}

// RateLimit represents a rate limit configuration
type RateLimit struct {
	MaxAttempts int           // Maximum attempts allowed
	Window      time.Duration // Time window for the attempts
	BlockTime   time.Duration // How long to block after exceeding limit
}

// Default rate limits for different actions
var defaultRateLimits = map[string]RateLimit{
	"login": {
		MaxAttempts: 5,
		Window:      15 * time.Minute,
		BlockTime:   30 * time.Minute,
	},
	"register": {
		MaxAttempts: 3,
		Window:      time.Hour,
		BlockTime:   2 * time.Hour,
	},
	"password_reset": {
		MaxAttempts: 3,
		Window:      time.Hour,
		BlockTime:   time.Hour,
	},
	"sms_send": {
		MaxAttempts: 5,
		Window:      time.Hour,
		BlockTime:   2 * time.Hour,
	},
	"email_verification": {
		MaxAttempts: 3,
		Window:      time.Hour,
		BlockTime:   time.Hour,
	},
}

// NewRateLimiter creates a new rate limiter instance
func NewRateLimiter(db *sql.DB) *RateLimiter {
	return &RateLimiter{db: db}
}

// IsAllowed checks if an action is allowed for the given identifier
func (r *RateLimiter) IsAllowed(identifier, action string) (bool, time.Duration, error) {
	limit, exists := defaultRateLimits[action]
	if !exists {
		// If no rate limit is defined, allow the action
		return true, 0, nil
	}

	// Check if currently blocked
	var blockedUntil *time.Time
	err := r.db.QueryRow(`
		SELECT blocked_until 
		FROM rate_limits 
		WHERE identifier = $1 AND action = $2
	`, identifier, action).Scan(&blockedUntil)

	if err != nil && err != sql.ErrNoRows {
		return false, 0, fmt.Errorf("failed to check rate limit: %w", err)
	}

	// If blocked, check if block period has expired
	if blockedUntil != nil && blockedUntil.After(time.Now()) {
		remaining := time.Until(*blockedUntil)
		return false, remaining, nil
	}

	// Check attempts in current window
	windowStart := time.Now().Add(-limit.Window)
	var attempts int
	var recordExists bool

	err = r.db.QueryRow(`
		SELECT attempts, EXISTS(SELECT 1 FROM rate_limits WHERE identifier = $1 AND action = $2)
		FROM rate_limits 
		WHERE identifier = $1 AND action = $2 AND window_start > $3
	`, identifier, action, windowStart).Scan(&attempts, &recordExists)

	if err != nil && err != sql.ErrNoRows {
		return false, 0, fmt.Errorf("failed to get attempt count: %w", err)
	}

	// If no record exists or window has expired, reset the counter
	if err == sql.ErrNoRows || !recordExists {
		attempts = 0
	}

	// Check if limit exceeded
	if attempts >= limit.MaxAttempts {
		// Block the identifier
		blockedUntil := time.Now().Add(limit.BlockTime)
		_, err = r.db.Exec(`
			INSERT INTO rate_limits (identifier, action, attempts, window_start, blocked_until)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (identifier, action)
			DO UPDATE SET 
				attempts = $3,
				window_start = $4,
				blocked_until = $5,
				updated_at = NOW()
		`, identifier, action, attempts+1, time.Now(), blockedUntil)

		if err != nil {
			return false, 0, fmt.Errorf("failed to update rate limit: %w", err)
		}

		return false, limit.BlockTime, nil
	}

	return true, 0, nil
}

// RecordAttempt records an attempt for the given identifier and action
func (r *RateLimiter) RecordAttempt(identifier, action string) error {
	limit, exists := defaultRateLimits[action]
	if !exists {
		// If no rate limit is defined, don't record anything
		return nil
	}

	windowStart := time.Now().Add(-limit.Window)
	
	// Get current attempts in window
	var attempts int
	err := r.db.QueryRow(`
		SELECT COALESCE(attempts, 0)
		FROM rate_limits 
		WHERE identifier = $1 AND action = $2 AND window_start > $3
	`, identifier, action, windowStart).Scan(&attempts)

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to get current attempts: %w", err)
	}

	// If no record or window expired, start fresh
	if err == sql.ErrNoRows {
		attempts = 0
	}

	// Increment attempts
	attempts++
	
	// Check if this attempt exceeds the limit
	var blockedUntil *time.Time
	if attempts >= limit.MaxAttempts {
		blockTime := time.Now().Add(limit.BlockTime)
		blockedUntil = &blockTime
	}

	// Update or insert the record
	_, err = r.db.Exec(`
		INSERT INTO rate_limits (identifier, action, attempts, window_start, blocked_until)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (identifier, action)
		DO UPDATE SET 
			attempts = CASE 
				WHEN rate_limits.window_start <= $4 THEN $3
				ELSE rate_limits.attempts + 1
			END,
			window_start = CASE 
				WHEN rate_limits.window_start <= $4 THEN $4
				ELSE rate_limits.window_start
			END,
			blocked_until = $5,
			updated_at = NOW()
	`, identifier, action, attempts, time.Now(), blockedUntil)

	return err
}

// ResetAttempts resets the attempt counter for a given identifier and action
func (r *RateLimiter) ResetAttempts(identifier, action string) error {
	_, err := r.db.Exec(`
		UPDATE rate_limits 
		SET attempts = 0, blocked_until = NULL, updated_at = NOW()
		WHERE identifier = $1 AND action = $2
	`, identifier, action)
	return err
}

// GetRemainingAttempts returns the number of remaining attempts for an identifier/action
func (r *RateLimiter) GetRemainingAttempts(identifier, action string) (int, error) {
	limit, exists := defaultRateLimits[action]
	if !exists {
		return 0, fmt.Errorf("no rate limit defined for action: %s", action)
	}

	var attempts int

	err := r.db.QueryRow(`
		SELECT COALESCE(attempts, 0)
		FROM rate_limits 
		WHERE identifier = $1 AND action = $2 AND window_start > $3
	`, identifier, action, time.Now().Add(-limit.Window)).Scan(&attempts)

	if err != nil && err != sql.ErrNoRows {
		return 0, fmt.Errorf("failed to get attempts: %w", err)
	}

	remaining := limit.MaxAttempts - attempts
	if remaining < 0 {
		remaining = 0
	}

	return remaining, nil
}

// CleanupExpiredRecords removes old rate limit records
func (r *RateLimiter) CleanupExpiredRecords() error {
	// Remove records older than 24 hours that are not currently blocking
	_, err := r.db.Exec(`
		DELETE FROM rate_limits 
		WHERE window_start < NOW() - INTERVAL '24 hours' 
		AND (blocked_until IS NULL OR blocked_until < NOW())
	`)
	return err
}

// GetBlockStatus returns the block status for an identifier/action
func (r *RateLimiter) GetBlockStatus(identifier, action string) (bool, time.Duration, error) {
	var blockedUntil *time.Time
	err := r.db.QueryRow(`
		SELECT blocked_until 
		FROM rate_limits 
		WHERE identifier = $1 AND action = $2
	`, identifier, action).Scan(&blockedUntil)

	if err != nil && err != sql.ErrNoRows {
		return false, 0, fmt.Errorf("failed to check block status: %w", err)
	}

	if err == sql.ErrNoRows || blockedUntil == nil {
		return false, 0, nil
	}

	if blockedUntil.After(time.Now()) {
		remaining := time.Until(*blockedUntil)
		return true, remaining, nil
	}

	return false, 0, nil
}

// UnblockIdentifier removes a block for a specific identifier/action
func (r *RateLimiter) UnblockIdentifier(identifier, action string) error {
	_, err := r.db.Exec(`
		UPDATE rate_limits 
		SET blocked_until = NULL, updated_at = NOW()
		WHERE identifier = $1 AND action = $2
	`, identifier, action)
	return err
}

// GetRateLimitInfo returns comprehensive rate limit information
type RateLimitInfo struct {
	Action          string        `json:"action"`
	MaxAttempts     int           `json:"max_attempts"`
	CurrentAttempts int           `json:"current_attempts"`
	RemainingAttempts int         `json:"remaining_attempts"`
	WindowDuration  time.Duration `json:"window_duration"`
	IsBlocked       bool          `json:"is_blocked"`
	BlockedUntil    *time.Time    `json:"blocked_until,omitempty"`
	TimeRemaining   time.Duration `json:"time_remaining,omitempty"`
}

// GetRateLimitInfo returns detailed rate limit information for an identifier/action
func (r *RateLimiter) GetRateLimitInfo(identifier, action string) (*RateLimitInfo, error) {
	limit, exists := defaultRateLimits[action]
	if !exists {
		return nil, fmt.Errorf("no rate limit defined for action: %s", action)
	}

	var attempts int
	var blockedUntil *time.Time

	err := r.db.QueryRow(`
		SELECT COALESCE(attempts, 0), blocked_until
		FROM rate_limits 
		WHERE identifier = $1 AND action = $2
	`, identifier, action).Scan(&attempts, &blockedUntil)

	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get rate limit info: %w", err)
	}

	// If no record or window expired, reset attempts
	if err == sql.ErrNoRows {
		attempts = 0
	}

	info := &RateLimitInfo{
		Action:            action,
		MaxAttempts:       limit.MaxAttempts,
		CurrentAttempts:   attempts,
		RemainingAttempts: limit.MaxAttempts - attempts,
		WindowDuration:    limit.Window,
		IsBlocked:         false,
	}

	// Check if blocked
	if blockedUntil != nil && blockedUntil.After(time.Now()) {
		info.IsBlocked = true
		info.BlockedUntil = blockedUntil
		info.TimeRemaining = time.Until(*blockedUntil)
	}

	if info.RemainingAttempts < 0 {
		info.RemainingAttempts = 0
	}

	return info, nil
}