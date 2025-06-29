package middleware

import (
	"net/http"
	"strings"

	"arvfinder-backend/database"
	"arvfinder-backend/services"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware creates an authentication middleware
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Authorization header required",
			})
			c.Abort()
			return
		}

		// Check if it's a Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		token := parts[1]

		// Initialize auth service
		db := database.GetDB()
		jwtSecret := "your-super-secret-jwt-key-change-in-production" // Should come from env
		authService := services.NewAuthService(db, jwtSecret)

		// Validate the token
		claims, err := authService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("tenant_id", claims.TenantID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("session_id", claims.SessionID)

		c.Next()
	}
}

// SecurityHeadersMiddleware adds security headers
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent XSS attacks
		c.Header("X-XSS-Protection", "1; mode=block")
		
		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")
		
		// Prevent clickjacking
		c.Header("X-Frame-Options", "DENY")
		
		// Force HTTPS (in production)
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		
		// Content Security Policy
		c.Header("Content-Security-Policy", "default-src 'self'")
		
		// Referrer Policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// Feature Policy
		c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=()")

		c.Next()
	}
}

// RateLimitMiddleware adds rate limiting
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// This would integrate with the rate limiter service
		// For now, just continue
		c.Next()
	}
}