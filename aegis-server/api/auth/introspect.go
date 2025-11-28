// Package auth provides HTTP REST API endpoints for authentication and token management.
// This file implements RFC 7662 OAuth 2.0 Token Introspection endpoint.
package auth

import (
	"log"
	"net/http"
	"nfcunha/aegis/domain/token"
	"strings"
	"github.com/gin-gonic/gin"
	"nfcunha/aegis/util/jwt"
)

// IntrospectTokenRequest represents the request body for token introspection endpoint.
// Follows RFC 7662 OAuth 2.0 Token Introspection specification.
type IntrospectTokenRequest struct {
	Token         string `json:"token" binding:"required"`
	TokenTypeHint string `json:"token_type_hint,omitempty"` // "access_token" or "refresh_token"
}

// IntrospectTokenResponse represents the response structure for token introspection.
// Follows RFC 7662 section 2.2 - Introspection Response.
//
// For active tokens, returns full metadata including OAuth2 standard claims.
// For inactive tokens, returns only {"active": false} per RFC 7662.
type IntrospectTokenResponse struct {
	// Active is REQUIRED. Boolean indicator of whether the token is currently active.
	Active bool `json:"active"`
	
	// The following fields are OPTIONAL and only included when Active is true:
	
	// Scope is a space-separated list of scopes associated with the token.
	Scope string `json:"scope,omitempty"`
	
	// ClientId is the identifier for the OAuth 2.0 client that requested the token.
	ClientId string `json:"client_id,omitempty"`
	
	// Username is a human-readable identifier for the resource owner (typically email).
	Username string `json:"username,omitempty"`
	
	// TokenType is the type of token (typically "Bearer").
	TokenType string `json:"token_type,omitempty"`
	
	// Exp is the Unix timestamp indicating when the token expires.
	Exp int64 `json:"exp,omitempty"`
	
	// Iat is the Unix timestamp indicating when the token was issued.
	Iat int64 `json:"iat,omitempty"`
	
	// Sub is the subject identifier (user ID).
	Sub string `json:"sub,omitempty"`
	
	// Iss is the issuer identifier (who issued the token).
	Iss string `json:"iss,omitempty"`
	
	// Extension fields (not part of RFC 7662 but useful for Aegis):
	
	// Roles contains the list of roles assigned to the user.
	Roles []string `json:"roles,omitempty"`
	
	// Permissions contains the list of permissions granted to the user.
	Permissions []string `json:"permissions,omitempty"`
}

// IntrospectToken is an HTTP handler that implements RFC 7662 OAuth 2.0 Token Introspection.
// It validates tokens and returns metadata in OAuth2-compliant format.
//
// Endpoint: POST /aegis/api/auth/introspect
//
// Request Body:
//   - token: The token to introspect (required)
//   - token_type_hint: Optional hint about the token type ("access_token" or "refresh_token")
//
// Response (200 OK):
//   - For active tokens: Returns active=true with full OAuth2 metadata
//   - For inactive tokens: Returns only {"active": false}
//
// The endpoint always returns 200 OK status per RFC 7662 section 2.2.
// This allows clients to distinguish between network errors and validation results.
//
// Standards Compliance:
//   - RFC 7662: OAuth 2.0 Token Introspection
//   - Compatible with OAuth2 API gateways and proxies
func IntrospectToken(c *gin.Context) {
	log.Println("POST /aegis/api/auth/introspect - Token introspection request received")
	
	var req IntrospectTokenRequest
	
	// Parse and validate request body
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Log token type hint if provided
	if req.TokenTypeHint != "" {
		log.Printf("Token type hint: %s", req.TokenTypeHint)
	}
	
	// Validate the token using JWT utility
	claims, err := jwt.ValidateToken(req.Token)
	
	// Handle validation errors - return inactive token response per RFC 7662
	if err != nil {
		log.Printf("Token introspection failed: %v", err)
		
		// RFC 7662 Section 2.2: Return minimal response for inactive tokens
		c.JSON(http.StatusOK, IntrospectTokenResponse{
			Active: false,
		})
		return
	}
	
	// Check if token is blacklisted (revoked)
	if token.GlobalBlacklist != nil && token.GlobalBlacklist.IsBlacklisted(claims.ID) {
		log.Printf("Token introspection failed: token has been revoked (JTI: %s)", claims.ID)
		
		// RFC 7662 Section 2.2: Return minimal response for inactive tokens
		c.JSON(http.StatusOK, IntrospectTokenResponse{
			Active: false,
		})
		return
	}
	
	// Token is active - return full OAuth2 metadata
	log.Printf("Token introspection successful for user: %s", claims.Subject)
	
	// Build scope string from roles and permissions
	// Format: "role:admin role:manager permission:read:users permission:write:users"
	scope := buildScopeString(claims.Roles, claims.Permissions)
	
	// Construct RFC 7662-compliant response
	response := IntrospectTokenResponse{
		Active:      true,
		Scope:       scope,
		ClientId:    "aegis-default-client", // TODO: Implement client management in Phase 3
		Username:    claims.Subject,
		TokenType:   "Bearer",
		Exp:         claims.ExpiresAt.Unix(),
		Iat:         claims.IssuedAt.Unix(),
		Sub:         claims.UserId,
		Iss:         claims.Issuer,
		Roles:       claims.Roles,
		Permissions: claims.Permissions,
	}
	
	c.JSON(http.StatusOK, response)
}

// buildScopeString constructs an OAuth2-compliant scope string from roles and permissions.
// Scope format: Space-separated list of scope identifiers.
//
// For Aegis, we prefix roles with "role:" and permissions are used as-is:
//   - Roles: "role:admin role:manager"
//   - Permissions: "read:users write:users"
//   - Combined: "role:admin role:manager read:users write:users"
//
// Parameters:
//   - roles: List of role names
//   - permissions: List of permission names
//
// Returns:
//   - Space-separated scope string
func buildScopeString(roles []string, permissions []string) string {
	var scopes []string
	
	// Add roles with "role:" prefix
	for _, role := range roles {
		if role != "" {
			scopes = append(scopes, "role:"+role)
		}
	}
	
	// Add permissions as-is (already in "resource:action" format)
	for _, permission := range permissions {
		if permission != "" {
			scopes = append(scopes, permission)
		}
	}
	
	// Join with spaces per OAuth2 specification
	return strings.Join(scopes, " ")
}
