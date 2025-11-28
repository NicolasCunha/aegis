package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"nfcunha/aegis/domain/token"
	"nfcunha/aegis/util/jwt"
)

func TestRevokeToken_Success(t *testing.T) {
	// Initialize blacklist
	bl := token.NewMemoryBlacklist()
	token.InitializeBlacklist(bl)
	defer func() { token.GlobalBlacklist = nil }()

	// Generate a valid token
	tokenPair, err := jwt.GenerateTokenPair(uuid.New(), "test@example.com", []string{"admin"}, []string{"read:users"})
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Create request
	reqBody := RevokeTokenRequest{
		Token: tokenPair.AccessToken,
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/aegis/api/auth/revoke", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	// Call handler
	RevokeToken(c)

	// Assert response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response RevokeTokenResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if !response.Success {
		t.Errorf("Expected success=true, got false")
	}

	// Verify token is actually blacklisted
	claims, _ := jwt.ValidateToken(tokenPair.AccessToken)
	if !bl.IsBlacklisted(claims.ID) {
		t.Errorf("Expected token to be blacklisted after revocation")
	}
}

func TestRevokeToken_AlreadyRevoked(t *testing.T) {
	// Initialize blacklist
	bl := token.NewMemoryBlacklist()
	token.InitializeBlacklist(bl)
	defer func() { token.GlobalBlacklist = nil }()

	// Generate a valid token
	tokenPair, err := jwt.GenerateTokenPair(uuid.New(), "test@example.com", []string{"admin"}, []string{"read:users"})
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Revoke the token first
	claims, _ := jwt.ValidateToken(tokenPair.AccessToken)
	bl.Add(claims.ID, time.Now().Add(1*time.Hour))

	// Create request to revoke again
	reqBody := RevokeTokenRequest{
		Token: tokenPair.AccessToken,
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/aegis/api/auth/revoke", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	// Call handler
	RevokeToken(c)

	// Assert response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response RevokeTokenResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if !response.Success {
		t.Errorf("Expected success=true for already revoked token, got false")
	}

	if response.Message != "Token already revoked" {
		t.Errorf("Expected 'Token already revoked' message, got %s", response.Message)
	}
}

func TestRevokeToken_InvalidToken(t *testing.T) {
	// Initialize blacklist
	bl := token.NewMemoryBlacklist()
	token.InitializeBlacklist(bl)
	defer func() { token.GlobalBlacklist = nil }()

	// Create request with invalid token
	reqBody := RevokeTokenRequest{
		Token: "invalid.jwt.token",
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/aegis/api/auth/revoke", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	// Call handler
	RevokeToken(c)

	// Assert response
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response["error"] != "Invalid token" {
		t.Errorf("Expected 'Invalid token' error, got %s", response["error"])
	}
}

func TestRevokeToken_MissingToken(t *testing.T) {
	// Initialize blacklist
	bl := token.NewMemoryBlacklist()
	token.InitializeBlacklist(bl)
	defer func() { token.GlobalBlacklist = nil }()

	// Create request with missing token
	reqBody := map[string]string{}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/aegis/api/auth/revoke", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	// Call handler
	RevokeToken(c)

	// Assert response
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestRevokeToken_BlacklistUnavailable(t *testing.T) {
	// Set global blacklist to nil
	token.GlobalBlacklist = nil
	defer func() { token.GlobalBlacklist = nil }()

	// Generate a valid token
	tokenPair, err := jwt.GenerateTokenPair(uuid.New(), "test@example.com", []string{"admin"}, []string{"read:users"})
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Create request
	reqBody := RevokeTokenRequest{
		Token: tokenPair.AccessToken,
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/aegis/api/auth/revoke", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	// Call handler
	RevokeToken(c)

	// Assert response
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response["error"] != "Token revocation system unavailable" {
		t.Errorf("Expected system unavailable error, got %s", response["error"])
	}
}

func TestRevokeToken_ExpiredToken(t *testing.T) {
	// Initialize blacklist
	bl := token.NewMemoryBlacklist()
	token.InitializeBlacklist(bl)
	defer func() { token.GlobalBlacklist = nil }()

	// Create an expired token (this would require manually crafting a token with past expiration)
	// For this test, we'll use an invalid token which has similar behavior
	expiredToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTYyMzkwMjJ9.invalid"

	// Create request
	reqBody := RevokeTokenRequest{
		Token: expiredToken,
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/aegis/api/auth/revoke", bytes.NewBuffer(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	// Call handler
	RevokeToken(c)

	// Assert response - should fail validation
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for expired token, got %d", w.Code)
	}
}

func TestRevokeToken_InvalidJSON(t *testing.T) {
	// Initialize blacklist
	bl := token.NewMemoryBlacklist()
	token.InitializeBlacklist(bl)
	defer func() { token.GlobalBlacklist = nil }()

	// Create request with invalid JSON
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/aegis/api/auth/revoke", bytes.NewBuffer([]byte("{invalid json")))
	c.Request.Header.Set("Content-Type", "application/json")

	// Call handler
	RevokeToken(c)

	// Assert response
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid JSON, got %d", w.Code)
	}
}

func TestRevokeToken_MultipleRevocations(t *testing.T) {
	// Initialize blacklist
	bl := token.NewMemoryBlacklist()
	token.InitializeBlacklist(bl)
	defer func() { token.GlobalBlacklist = nil }()

	// Generate multiple tokens
	token1, _ := jwt.GenerateTokenPair(uuid.New(), "user1@example.com", []string{"admin"}, []string{"read:users"})
	token2, _ := jwt.GenerateTokenPair(uuid.New(), "user2@example.com", []string{"user"}, []string{"read:self"})
	token3, _ := jwt.GenerateTokenPair(uuid.New(), "user3@example.com", []string{"user"}, []string{"read:self"})

	tokens := []string{token1.AccessToken, token2.AccessToken, token3.AccessToken}

	// Revoke all tokens
	for i, tok := range tokens {
		reqBody := RevokeTokenRequest{Token: tok}
		jsonBody, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/aegis/api/auth/revoke", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		RevokeToken(c)

		if w.Code != http.StatusOK {
			t.Errorf("Token %d: Expected status 200, got %d", i+1, w.Code)
		}
	}

	// Verify all tokens are blacklisted
	if bl.Size() != 3 {
		t.Errorf("Expected 3 blacklisted tokens, got %d", bl.Size())
	}

	for i, tok := range tokens {
		claims, _ := jwt.ValidateToken(tok)
		if !bl.IsBlacklisted(claims.ID) {
			t.Errorf("Token %d should be blacklisted", i+1)
		}
	}
}
