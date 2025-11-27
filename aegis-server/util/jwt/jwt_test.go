package jwt

import (
	"os"
	"strings"
	"testing"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// TestGenerateTokenPair tests successful token pair generation
func TestGenerateTokenPair(t *testing.T) {
	userId := uuid.New()
	subject := "test@example.com"
	roles := []string{"admin", "user"}
	permissions := []string{"read", "write"}
	
	tokenPair, err := GenerateTokenPair(userId, subject, roles, permissions)
	
	if err != nil {
		t.Fatalf("GenerateTokenPair should not return error: %v", err)
	}
	if tokenPair == nil {
		t.Fatal("TokenPair should not be nil")
	}
	if tokenPair.AccessToken == "" {
		t.Error("AccessToken should not be empty")
	}
	if tokenPair.RefreshToken == "" {
		t.Error("RefreshToken should not be empty")
	}
	if tokenPair.ExpiresAt.IsZero() {
		t.Error("ExpiresAt should not be zero")
	}
	if tokenPair.RefreshExpiresAt.IsZero() {
		t.Error("RefreshExpiresAt should not be zero")
	}
	
	// Verify refresh token expires after access token
	if !tokenPair.RefreshExpiresAt.After(tokenPair.ExpiresAt) {
		t.Error("RefreshToken should expire after AccessToken")
	}
}

// TestGenerateTokenPair_TokensAreDifferent tests that access and refresh tokens are different
func TestGenerateTokenPair_TokensAreDifferent(t *testing.T) {
	userId := uuid.New()
	subject := "test@example.com"
	roles := []string{"admin"}
	permissions := []string{"read"}
	
	tokenPair, err := GenerateTokenPair(userId, subject, roles, permissions)
	
	if err != nil {
		t.Fatalf("GenerateTokenPair should not return error: %v", err)
	}
	if tokenPair.AccessToken == tokenPair.RefreshToken {
		t.Error("AccessToken and RefreshToken should be different")
	}
}

// TestGenerateTokenPair_EmptyRolesAndPermissions tests token generation with empty roles/permissions
func TestGenerateTokenPair_EmptyRolesAndPermissions(t *testing.T) {
	userId := uuid.New()
	subject := "test@example.com"
	
	tokenPair, err := GenerateTokenPair(userId, subject, []string{}, []string{})
	
	if err != nil {
		t.Fatalf("GenerateTokenPair should not return error: %v", err)
	}
	if tokenPair == nil {
		t.Fatal("TokenPair should not be nil")
	}
}

// TestValidateToken_ValidAccessToken tests validation of a valid access token
func TestValidateToken_ValidAccessToken(t *testing.T) {
	userId := uuid.New()
	subject := "test@example.com"
	roles := []string{"admin", "user"}
	permissions := []string{"read", "write"}
	
	tokenPair, err := GenerateTokenPair(userId, subject, roles, permissions)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	
	claims, err := ValidateToken(tokenPair.AccessToken)
	
	if err != nil {
		t.Fatalf("ValidateToken should not return error: %v", err)
	}
	if claims.UserId != userId.String() {
		t.Errorf("Expected UserId %s, got %s", userId.String(), claims.UserId)
	}
	if claims.Subject != subject {
		t.Errorf("Expected Subject %s, got %s", subject, claims.Subject)
	}
	if len(claims.Roles) != len(roles) {
		t.Errorf("Expected %d roles, got %d", len(roles), len(claims.Roles))
	}
	if len(claims.Permissions) != len(permissions) {
		t.Errorf("Expected %d permissions, got %d", len(permissions), len(claims.Permissions))
	}
	if claims.TokenType != "access" {
		t.Errorf("Expected TokenType 'access', got '%s'", claims.TokenType)
	}
}

// TestValidateToken_ValidRefreshToken tests validation of a valid refresh token
func TestValidateToken_ValidRefreshToken(t *testing.T) {
	userId := uuid.New()
	subject := "test@example.com"
	roles := []string{"admin"}
	permissions := []string{"read"}
	
	tokenPair, err := GenerateTokenPair(userId, subject, roles, permissions)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	
	claims, err := ValidateToken(tokenPair.RefreshToken)
	
	if err != nil {
		t.Fatalf("ValidateToken should not return error: %v", err)
	}
	if claims.TokenType != "refresh" {
		t.Errorf("Expected TokenType 'refresh', got '%s'", claims.TokenType)
	}
}

// TestValidateToken_InvalidToken tests validation with an invalid token
func TestValidateToken_InvalidToken(t *testing.T) {
	invalidToken := "invalid.token.here"
	
	_, err := ValidateToken(invalidToken)
	
	if err == nil {
		t.Error("ValidateToken should return error for invalid token")
	}
}

// TestValidateToken_EmptyToken tests validation with empty token
func TestValidateToken_EmptyToken(t *testing.T) {
	_, err := ValidateToken("")
	
	if err == nil {
		t.Error("ValidateToken should return error for empty token")
	}
}

// TestValidateToken_MalformedToken tests validation with malformed token
func TestValidateToken_MalformedToken(t *testing.T) {
	malformedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.malformed"
	
	_, err := ValidateToken(malformedToken)
	
	if err == nil {
		t.Error("ValidateToken should return error for malformed token")
	}
}

// TestValidateToken_TamperedToken tests validation with tampered token
func TestValidateToken_TamperedToken(t *testing.T) {
	userId := uuid.New()
	subject := "test@example.com"
	
	tokenPair, err := GenerateTokenPair(userId, subject, []string{}, []string{})
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	
	// Tamper with the token by changing a character
	tamperedToken := tokenPair.AccessToken[:len(tokenPair.AccessToken)-5] + "XXXXX"
	
	_, err = ValidateToken(tamperedToken)
	
	if err == nil {
		t.Error("ValidateToken should return error for tampered token")
	}
}

// TestValidateToken_WrongSigningKey tests validation with different signing key
func TestValidateToken_WrongSigningKey(t *testing.T) {
	userId := uuid.New()
	subject := "test@example.com"
	
	// Generate token with current key
	tokenPair, err := GenerateTokenPair(userId, subject, []string{}, []string{})
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	
	// Change the signing key
	originalSecret := JWT_SECRET
	JWT_SECRET = "different_secret_key_for_testing"
	defer func() { JWT_SECRET = originalSecret }()
	
	// Try to validate with different key
	_, err = ValidateToken(tokenPair.AccessToken)
	
	if err == nil {
		t.Error("ValidateToken should return error when signing key is different")
	}
}

// TestValidateRefreshToken_ValidRefreshToken tests successful refresh token validation
func TestValidateRefreshToken_ValidRefreshToken(t *testing.T) {
	userId := uuid.New()
	subject := "test@example.com"
	
	tokenPair, err := GenerateTokenPair(userId, subject, []string{}, []string{})
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	
	claims, err := ValidateRefreshToken(tokenPair.RefreshToken)
	
	if err != nil {
		t.Fatalf("ValidateRefreshToken should not return error: %v", err)
	}
	if claims.TokenType != "refresh" {
		t.Errorf("Expected TokenType 'refresh', got '%s'", claims.TokenType)
	}
}

// TestValidateRefreshToken_AccessTokenRejected tests that access tokens are rejected
func TestValidateRefreshToken_AccessTokenRejected(t *testing.T) {
	userId := uuid.New()
	subject := "test@example.com"
	
	tokenPair, err := GenerateTokenPair(userId, subject, []string{}, []string{})
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	
	_, err = ValidateRefreshToken(tokenPair.AccessToken)
	
	if err == nil {
		t.Error("ValidateRefreshToken should return error for access token")
	}
	if !strings.Contains(err.Error(), "not a refresh token") {
		t.Errorf("Error should mention 'not a refresh token', got: %v", err)
	}
}

// TestValidateRefreshToken_InvalidToken tests refresh token validation with invalid token
func TestValidateRefreshToken_InvalidToken(t *testing.T) {
	_, err := ValidateRefreshToken("invalid.token")
	
	if err == nil {
		t.Error("ValidateRefreshToken should return error for invalid token")
	}
}

// TestGetTokenExpiration_DefaultValue tests default token expiration
func TestGetTokenExpiration_DefaultValue(t *testing.T) {
	os.Unsetenv("AEGIS_JWT_EXP_TIME")
	
	duration := getTokenExpiration()
	
	expectedDuration := 1440 * time.Minute // 24 hours
	if duration != expectedDuration {
		t.Errorf("Expected default expiration %v, got %v", expectedDuration, duration)
	}
}

// TestGetTokenExpiration_CustomValue tests custom token expiration from env var
func TestGetTokenExpiration_CustomValue(t *testing.T) {
	os.Setenv("AEGIS_JWT_EXP_TIME", "60")
	defer os.Unsetenv("AEGIS_JWT_EXP_TIME")
	
	duration := getTokenExpiration()
	
	expectedDuration := 60 * time.Minute
	if duration != expectedDuration {
		t.Errorf("Expected custom expiration %v, got %v", expectedDuration, duration)
	}
}

// TestGetTokenExpiration_InvalidValue tests invalid env var falls back to default
func TestGetTokenExpiration_InvalidValue(t *testing.T) {
	os.Setenv("AEGIS_JWT_EXP_TIME", "invalid")
	defer os.Unsetenv("AEGIS_JWT_EXP_TIME")
	
	duration := getTokenExpiration()
	
	expectedDuration := 1440 * time.Minute // Should fall back to default
	if duration != expectedDuration {
		t.Errorf("Expected default expiration %v, got %v", expectedDuration, duration)
	}
}

// TestGetTokenExpiration_NegativeValue tests negative value falls back to default
func TestGetTokenExpiration_NegativeValue(t *testing.T) {
	os.Setenv("AEGIS_JWT_EXP_TIME", "-10")
	defer os.Unsetenv("AEGIS_JWT_EXP_TIME")
	
	duration := getTokenExpiration()
	
	expectedDuration := 1440 * time.Minute // Should fall back to default
	if duration != expectedDuration {
		t.Errorf("Expected default expiration %v, got %v", expectedDuration, duration)
	}
}

// TestGetTokenExpiration_ZeroValue tests zero value falls back to default
func TestGetTokenExpiration_ZeroValue(t *testing.T) {
	os.Setenv("AEGIS_JWT_EXP_TIME", "0")
	defer os.Unsetenv("AEGIS_JWT_EXP_TIME")
	
	duration := getTokenExpiration()
	
	expectedDuration := 1440 * time.Minute // Should fall back to default
	if duration != expectedDuration {
		t.Errorf("Expected default expiration %v, got %v", expectedDuration, duration)
	}
}

// TestGetJwtSecret_CustomValue tests custom JWT secret from env var
func TestGetJwtSecret_CustomValue(t *testing.T) {
	customSecret := "my_custom_jwt_secret_256bit_key"
	os.Setenv("AEGIS_JWT_SECRET", customSecret)
	defer os.Unsetenv("AEGIS_JWT_SECRET")
	
	secret := getJwtSecret()
	
	if secret != customSecret {
		t.Errorf("Expected custom secret '%s', got '%s'", customSecret, secret)
	}
}

// TestGetJwtSecret_GeneratedValue tests generated JWT secret when env var not set
func TestGetJwtSecret_GeneratedValue(t *testing.T) {
	os.Unsetenv("AEGIS_JWT_SECRET")
	
	secret := getJwtSecret()
	
	if secret == "" {
		t.Error("Generated secret should not be empty")
	}
	// Should be 64 hex characters (32 bytes = 256 bits)
	if len(secret) != 64 {
		t.Errorf("Generated secret should be 64 characters, got %d", len(secret))
	}
}

// TestTokenExpiration tests that tokens expire after the configured time
func TestTokenExpiration(t *testing.T) {
	userId := uuid.New()
	subject := "test@example.com"
	
	// Set very short expiration for testing
	originalExpiration := TOKEN_EXPIRATION
	TOKEN_EXPIRATION = 1 * time.Millisecond
	defer func() { TOKEN_EXPIRATION = originalExpiration }()
	
	tokenPair, err := GenerateTokenPair(userId, subject, []string{}, []string{})
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	
	// Wait for token to expire
	time.Sleep(10 * time.Millisecond)
	
	_, err = ValidateToken(tokenPair.AccessToken)
	
	if err == nil {
		t.Error("ValidateToken should return error for expired token")
	}
}

// TestTokenClaims_RolesAndPermissions tests that roles and permissions are correctly embedded
func TestTokenClaims_RolesAndPermissions(t *testing.T) {
	userId := uuid.New()
	subject := "test@example.com"
	roles := []string{"admin", "moderator", "user"}
	permissions := []string{"read", "write", "delete", "admin:all"}
	
	tokenPair, err := GenerateTokenPair(userId, subject, roles, permissions)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	
	claims, err := ValidateToken(tokenPair.AccessToken)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}
	
	// Verify all roles are present
	for _, role := range roles {
		found := false
		for _, claimRole := range claims.Roles {
			if claimRole == role {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Role '%s' not found in token claims", role)
		}
	}
	
	// Verify all permissions are present
	for _, permission := range permissions {
		found := false
		for _, claimPermission := range claims.Permissions {
			if claimPermission == permission {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Permission '%s' not found in token claims", permission)
		}
	}
}

// TestTokenClaims_RegisteredClaims tests standard JWT claims
func TestTokenClaims_RegisteredClaims(t *testing.T) {
	userId := uuid.New()
	subject := "test@example.com"
	
	tokenPair, err := GenerateTokenPair(userId, subject, []string{}, []string{})
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	
	claims, err := ValidateToken(tokenPair.AccessToken)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}
	
	// Verify issuer
	if claims.Issuer != "aegis" {
		t.Errorf("Expected issuer 'aegis', got '%s'", claims.Issuer)
	}
	
	// Verify IssuedAt is set and in the past
	if claims.IssuedAt == nil {
		t.Error("IssuedAt should not be nil")
	} else if claims.IssuedAt.Time.After(time.Now()) {
		t.Error("IssuedAt should be in the past")
	}
	
	// Verify ExpiresAt is set and in the future
	if claims.ExpiresAt == nil {
		t.Error("ExpiresAt should not be nil")
	} else if claims.ExpiresAt.Time.Before(time.Now()) {
		t.Error("ExpiresAt should be in the future")
	}
}

// TestGenerateTokenWithType_InvalidSigningMethod tests handling of wrong signing method
func TestGenerateTokenWithType_InvalidSigningMethod(t *testing.T) {
	userId := uuid.New()
	subject := "test@example.com"
	
	// Generate a token with RS256 (RSA) instead of HS256 (HMAC)
	claims := &TokenClaims{
		UserId:      userId.String(),
		Subject:     subject,
		Roles:       []string{},
		Permissions: []string{},
		TokenType:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "aegis",
		},
	}
	
	// Create token with wrong signing method (this will be caught during validation)
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	tokenString, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	
	// Try to validate - should fail due to wrong signing method
	_, err := ValidateToken(tokenString)
	
	if err == nil {
		t.Error("ValidateToken should return error for wrong signing method")
	}
}
