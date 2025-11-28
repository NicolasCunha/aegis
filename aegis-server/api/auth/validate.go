// Package auth provides HTTP REST API endpoints for authentication and token management.
// Includes token validation, introspection, and revocation endpoints.
package auth

import (
	"log"
	"net/http"
	"nfcunha/aegis/domain/token"
	"strings"
	"time"
	"github.com/gin-gonic/gin"
	"nfcunha/aegis/util/jwt"
)

// ValidateTokenRequest represents the request body for token validation endpoint.
// It contains the JWT token string that needs to be validated.
type ValidateTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

// ValidateTokenResponse represents the response structure for token validation.
// When valid is true, it includes user claims and expiration information.
// When valid is false, it includes an error message describing why validation failed.
type ValidateTokenResponse struct {
	Valid     bool      `json:"valid"`
	User      *UserInfo `json:"user,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	Error     string    `json:"error,omitempty"`
}

// UserInfo represents the user claims extracted from a validated JWT token.
// It includes user identity, roles, and permissions for authorization purposes.
type UserInfo struct {
	ID          string   `json:"id"`
	Subject     string   `json:"subject"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}

// ValidateToken is an HTTP handler that validates JWT tokens and returns user claims.
// It accepts both access and refresh tokens and provides detailed validation results.
//
// Endpoint: POST /aegis/api/auth/validate
//
// Request Body:
//   - token: The JWT token string to validate (required)
//
// Response (200 OK):
//   - For valid tokens: Returns valid=true with user claims and expiration
//   - For invalid tokens: Returns valid=false with error description
//
// The endpoint always returns 200 OK status, with the validation result in the response body.
// This allows clients to distinguish between network errors and validation failures.
func ValidateToken(c *gin.Context) {
	log.Println("POST /aegis/api/auth/validate - Token validation request received")

	var req ValidateTokenRequest
	
	// Parse and validate request body
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the token using JWT utility
	claims, err := jwt.ValidateToken(req.Token)
	
	// Handle validation errors - return 200 with valid=false for invalid tokens
	if err != nil {
		log.Printf("Token validation failed: %v", err)
		
		// Determine error type for more descriptive messages
		errorMessage := determineValidationError(err)
		
		c.JSON(http.StatusOK, ValidateTokenResponse{
			Valid: false,
			Error: errorMessage,
		})
		return
	}

	// Check if token is blacklisted (revoked)
	if token.GlobalBlacklist != nil && token.GlobalBlacklist.IsBlacklisted(claims.ID) {
		log.Printf("Token is blacklisted (revoked): JTI=%s, User=%s", claims.ID, claims.Subject)
		c.JSON(http.StatusOK, ValidateTokenResponse{
			Valid: false,
			Error: "token revoked",
		})
		return
	}

	// Token is valid - return user claims and expiration
	log.Printf("Token validated successfully for user: %s", claims.Subject)
	
	// Extract expiration time from claims
	expiresAt := claims.ExpiresAt.Time
	
	c.JSON(http.StatusOK, ValidateTokenResponse{
		Valid: true,
		User: &UserInfo{
			ID:          claims.UserId,
			Subject:     claims.Subject,
			Roles:       claims.Roles,
			Permissions: claims.Permissions,
		},
		ExpiresAt: &expiresAt,
	})
}

// determineValidationError analyzes the validation error and returns a user-friendly message.
// This helps clients understand why token validation failed without exposing sensitive details.
//
// Parameters:
//   - err: The error returned from JWT validation
//
// Returns:
//   - A descriptive error message string
func determineValidationError(err error) string {
	errMsg := err.Error()
	
	// Check for common JWT validation errors using strings package
	switch {
	case strings.Contains(errMsg, "expired"):
		return "token expired"
	case strings.Contains(errMsg, "signature"):
		return "invalid signature"
	case strings.Contains(errMsg, "malformed"):
		return "malformed token"
	case strings.Contains(errMsg, "unexpected signing method"):
		return "invalid signing method"
	default:
		return "invalid token"
	}
}

// RegisterApi registers all auth-related HTTP routes with the Gin router.
// Includes token validation, introspection, and revocation endpoints per OAuth2/OIDC standards.
//
// These are public endpoints for client applications to validate tokens issued by Aegis.
// Client applications should implement their own authentication middleware using these endpoints.
//
// Public endpoints (under /aegis context path):
//   - POST /api/auth/validate - Validates a JWT token and returns user claims
//   - POST /api/auth/introspect - OAuth2-compliant token introspection (RFC 7662)
//   - POST /api/auth/revoke - Revokes a JWT token by adding it to the blacklist
//
// Parameters:
//   - router: The Gin RouterGroup to register routes with (already under /aegis)
func RegisterApi(router gin.IRouter) {
	auth := router.Group("/api/auth")
	{
		auth.POST("/validate", ValidateToken)
		auth.POST("/introspect", IntrospectToken)
		auth.POST("/revoke", RevokeToken)
	}
}
