// Package auth provides HTTP REST API endpoints for authentication and token management.
// This file implements token revocation functionality.
package auth

import (
	"log"
	"net/http"
	"nfcunha/aegis/domain/token"
	"time"

	"github.com/gin-gonic/gin"
	"nfcunha/aegis/util/jwt"
)

// RevokeTokenRequest represents the request structure for token revocation.
type RevokeTokenRequest struct {
	// Token is the JWT token to revoke (required)
	Token string `json:"token" binding:"required"`
}

// RevokeTokenResponse represents the response structure for token revocation.
type RevokeTokenResponse struct {
	// Success indicates whether the revocation was successful
	Success bool `json:"success"`
	
	// Message provides details about the operation
	Message string `json:"message"`
}

// RevokeToken is an HTTP handler that revokes a JWT token by adding it to the blacklist.
//
// Endpoint: POST /aegis/api/auth/revoke
//
// Request Body:
//   - token: The JWT token to revoke (required)
//
// Response:
//   - 200 OK: Token successfully revoked
//   - 400 Bad Request: Invalid request or token validation failed
//   - 500 Internal Server Error: Blacklist system unavailable
//
// The revoked token will be blacklisted until its natural expiration time.
// Subsequent validation or introspection requests for this token will return inactive/invalid.
func RevokeToken(c *gin.Context) {
	log.Println("POST /aegis/api/auth/revoke - Token revocation request received")
	
	var req RevokeTokenRequest
	
	// Parse and validate request body
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Check if blacklist is available
	if token.GlobalBlacklist == nil {
		log.Println("Token revocation failed: blacklist system not initialized")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Token revocation system unavailable",
		})
		return
	}
	
	// Validate the token first to ensure it's valid before revoking
	claims, err := jwt.ValidateToken(req.Token)
	if err != nil {
		log.Printf("Token revocation failed: invalid token - %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid token",
		})
		return
	}
	
	// Check if token is already blacklisted
	if token.GlobalBlacklist.IsBlacklisted(claims.ID) {
		log.Printf("Token already revoked (JTI: %s)", claims.ID)
		c.JSON(http.StatusOK, RevokeTokenResponse{
			Success: true,
			Message: "Token already revoked",
		})
		return
	}
	
	// Add token to blacklist
	expiresAt := time.Unix(claims.ExpiresAt.Unix(), 0)
	token.GlobalBlacklist.Add(claims.ID, expiresAt)
	
	log.Printf("Token revoked successfully (JTI: %s, User: %s)", claims.ID, claims.Subject)
	
	c.JSON(http.StatusOK, RevokeTokenResponse{
		Success: true,
		Message: "Token revoked successfully",
	})
}
