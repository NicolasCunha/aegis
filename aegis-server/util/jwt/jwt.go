// Package jwt provides utilities for generating and validating JSON Web Tokens (JWT)
// with HMAC-SHA256 signing. Tokens include user identity, roles, and permissions.
package jwt

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"os"
	"strconv"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var JWT_SECRET = getJwtSecret()
var TOKEN_EXPIRATION = getTokenExpiration()
const REFRESH_TOKEN_EXTRA_TIME = 1 * time.Minute

// TokenClaims represents the JWT claims structure containing user identity and authorization data.
// It embeds jwt.RegisteredClaims for standard JWT fields like expiration and issuer.
type TokenClaims struct {
	UserId      string   `json:"user_id"`
	Subject     string   `json:"subject"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
	TokenType   string   `json:"token_type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

// TokenOutput represents the result of token generation, containing the signed token
// string and its expiration timestamp.
type TokenOutput struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// TokenPair represents both access and refresh tokens returned during authentication.
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at"`
}

// GenerateTokenPair creates both an access token and a refresh token.
// The refresh token expires 1 minute after the access token to allow for token refresh.
//
// Parameters:
//   - userId: Unique identifier for the user
//   - subject: User's subject (typically email or username)
//   - roles: List of roles assigned to the user
//   - permissions: List of permissions granted to the user
//
// Returns:
//   - TokenPair containing both access and refresh tokens with their expiration times
//   - Error if token signing fails
func GenerateTokenPair(userId uuid.UUID, subject string, roles []string, permissions []string) (*TokenPair, error) {
	accessExpiration := time.Now().Add(TOKEN_EXPIRATION)
	refreshExpiration := time.Now().Add(TOKEN_EXPIRATION + REFRESH_TOKEN_EXTRA_TIME)

	// Generate access token
	accessToken, err := generateTokenWithType(userId, subject, roles, permissions, "access", TOKEN_EXPIRATION)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := generateTokenWithType(userId, subject, roles, permissions, "refresh", TOKEN_EXPIRATION + REFRESH_TOKEN_EXTRA_TIME)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:      accessToken.Token,
		RefreshToken:     refreshToken.Token,
		ExpiresAt:        accessExpiration,
		RefreshExpiresAt: refreshExpiration,
	}, nil
}

// generateTokenWithType creates a JWT token with a specific type (access or refresh).
func generateTokenWithType(userId uuid.UUID, subject string, roles []string, permissions []string, tokenType string, expiration time.Duration) (*TokenOutput, error) {
	expirationTime := time.Now().Add(expiration)

	claims := &TokenClaims{
		UserId:      userId.String(),
		Subject:     subject,
		Roles:       roles,
		Permissions: permissions,
		TokenType:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "aegis",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(JWT_SECRET))
	if err != nil {
		return nil, err
	}

	return &TokenOutput{
		Token:     tokenString,
		ExpiresAt: expirationTime,
	}, nil
}

// ValidateToken parses and validates a JWT token string, verifying its signature and expiration.
// The token must be signed with HMAC-SHA256 using the configured secret.
//
// Parameters:
//   - tokenString: The JWT token string to validate
//
// Returns:
//   - TokenClaims containing the extracted user information and authorization data
//   - Error if the token is invalid, expired, or has an unexpected signing method
func ValidateToken(tokenString string) (*TokenClaims, error) {
	claims := &TokenClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(JWT_SECRET), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// ValidateRefreshToken validates a refresh token and ensures it's of type "refresh".
//
// Parameters:
//   - tokenString: The refresh token string to validate
//
// Returns:
//   - TokenClaims containing the extracted user information
//   - Error if the token is invalid, expired, or not a refresh token
func ValidateRefreshToken(tokenString string) (*TokenClaims, error) {
	claims, err := ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != "refresh" {
		return nil, errors.New("token is not a refresh token")
	}

	return claims, nil
}

// getTokenExpiration retrieves the token expiration duration from AEGIS_JWT_EXP_TIME environment variable.
// The value should be in minutes. Defaults to 1440 minutes (24 hours) if not set.
//
// Returns:
//   - Token expiration duration
func getTokenExpiration() time.Duration {
	const JWT_EXP_TIME_ENV = "AEGIS_JWT_EXP_TIME"
	const DEFAULT_EXPIRATION_MINUTES = 1440 // 24 hours
	
	if expStr := os.Getenv(JWT_EXP_TIME_ENV); expStr != "" {
		if minutes, err := strconv.Atoi(expStr); err == nil && minutes > 0 {
			log.Printf("Using token expiration: %d minutes", minutes)
			return time.Duration(minutes) * time.Minute
		}
		log.Printf("Warning: invalid %s value '%s', using default %d minutes", JWT_EXP_TIME_ENV, expStr, DEFAULT_EXPIRATION_MINUTES)
	}
	
	log.Printf("Using default token expiration: %d minutes (24 hours)", DEFAULT_EXPIRATION_MINUTES)
	return time.Duration(DEFAULT_EXPIRATION_MINUTES) * time.Minute
}

// getJwtSecret retrieves the JWT signing secret from the AEGIS_JWT_SECRET environment variable.
// If not set, it generates a cryptographically secure 256-bit random secret suitable for HMAC-SHA256.
// A warning is logged when using a randomly generated secret.
//
// Returns:
//   - The JWT secret as a hex-encoded string
func getJwtSecret() string {
	const JWT_SECRET_ENV = "AEGIS_JWT_SECRET"
	if secret := os.Getenv(JWT_SECRET_ENV); secret != "" {
		return secret
	}
	// Generate 256-bit (32 bytes) secret for HMAC-SHA256
	secretBytes := make([]byte, 32)
	_, err := rand.Read(secretBytes)
	if err != nil {
		log.Fatal("Failed to generate JWT secret:", err)
	}
	generatedSecret := hex.EncodeToString(secretBytes)
	log.Printf("Warning: using randomly generated JWT secret, consider setting the environment variable '%s'\n", JWT_SECRET_ENV)
	log.Printf("Generated JWT secret: %s\n", generatedSecret)
	return generatedSecret
}
