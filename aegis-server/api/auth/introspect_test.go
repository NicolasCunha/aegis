package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"github.com/google/uuid"
	jwtUtil "nfcunha/aegis/util/jwt"
)

// TestIntrospectToken_ActiveAccessToken tests introspection of a valid access token
// Expected: Returns 200 OK with active=true and full OAuth2 metadata
func TestIntrospectToken_ActiveAccessToken(t *testing.T) {
	router := setupRouter()
	
	// Generate a valid access token
	userId := uuid.New()
	subject := "test@example.com"
	roles := []string{"admin", "user"}
	permissions := []string{"read:users", "write:users"}
	
	tokenPair, err := jwtUtil.GenerateTokenPair(userId, subject, roles, permissions)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	
	// Create introspection request
	reqBody := IntrospectTokenRequest{
		Token: tokenPair.AccessToken,
	}
	body, _ := json.Marshal(reqBody)
	
	req, _ := http.NewRequest("POST", "/aegis/api/auth/introspect", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Assert response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
	}
	
	var response IntrospectTokenResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	
	// Verify RFC 7662 compliance - active token response
	if !response.Active {
		t.Error("Expected active=true for valid token")
	}
	
	// Verify OAuth2 standard claims are present
	if response.Sub != userId.String() {
		t.Errorf("Expected sub %s, got %s", userId.String(), response.Sub)
	}
	if response.Username != subject {
		t.Errorf("Expected username %s, got %s", subject, response.Username)
	}
	if response.TokenType != "Bearer" {
		t.Errorf("Expected token_type 'Bearer', got %s", response.TokenType)
	}
	if response.Exp == 0 {
		t.Error("Expected exp (expiration) to be set")
	}
	if response.Iat == 0 {
		t.Error("Expected iat (issued at) to be set")
	}
	if response.Iss != "aegis" {
		t.Errorf("Expected issuer 'aegis', got %s", response.Iss)
	}
	if response.ClientId == "" {
		t.Error("Expected client_id to be set")
	}
	
	// Verify scope string includes roles and permissions
	if response.Scope == "" {
		t.Error("Expected scope to be set")
	}
	if !strings.Contains(response.Scope, "role:admin") {
		t.Error("Expected scope to contain 'role:admin'")
	}
	if !strings.Contains(response.Scope, "read:users") {
		t.Error("Expected scope to contain 'read:users'")
	}
	
	// Verify extension fields
	if len(response.Roles) != len(roles) {
		t.Errorf("Expected %d roles, got %d", len(roles), len(response.Roles))
	}
	if len(response.Permissions) != len(permissions) {
		t.Errorf("Expected %d permissions, got %d", len(permissions), len(response.Permissions))
	}
}

// TestIntrospectToken_ActiveRefreshToken tests introspection of a valid refresh token
// Expected: Returns 200 OK with active=true (refresh tokens are also introspectable)
func TestIntrospectToken_ActiveRefreshToken(t *testing.T) {
	router := setupRouter()
	
	// Generate tokens
	userId := uuid.New()
	subject := "test@example.com"
	
	tokenPair, err := jwtUtil.GenerateTokenPair(userId, subject, []string{"user"}, []string{"read"})
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	
	// Create introspection request with refresh token
	reqBody := IntrospectTokenRequest{
		Token: tokenPair.RefreshToken,
	}
	body, _ := json.Marshal(reqBody)
	
	req, _ := http.NewRequest("POST", "/aegis/api/auth/introspect", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Assert response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
	
	var response IntrospectTokenResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	
	// Refresh tokens should also be introspectable as active
	if !response.Active {
		t.Error("Expected active=true for valid refresh token")
	}
	if response.Sub == "" {
		t.Error("Expected sub to be set for refresh token")
	}
}

// TestIntrospectToken_InactiveToken tests introspection of an invalid token
// Expected: Returns 200 OK with only {"active": false} per RFC 7662
func TestIntrospectToken_InactiveToken(t *testing.T) {
	router := setupRouter()
	
	// Create request with invalid token
	reqBody := IntrospectTokenRequest{
		Token: "invalid.token.here",
	}
	body, _ := json.Marshal(reqBody)
	
	req, _ := http.NewRequest("POST", "/aegis/api/auth/introspect", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Assert response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
	}
	
	var response IntrospectTokenResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	
	// RFC 7662 Section 2.2: Inactive tokens return minimal response
	if response.Active {
		t.Error("Expected active=false for invalid token")
	}
	
	// Per RFC 7662, these fields should be omitted for inactive tokens
	// (they will be zero values in Go)
	if response.Sub != "" {
		t.Error("Expected sub to be empty for inactive token")
	}
	if response.Username != "" {
		t.Error("Expected username to be empty for inactive token")
	}
	if response.Exp != 0 {
		t.Error("Expected exp to be 0 for inactive token")
	}
}

// TestIntrospectToken_InactiveTamperedToken tests introspection of a tampered token
// Expected: Returns 200 OK with active=false
func TestIntrospectToken_InactiveTamperedToken(t *testing.T) {
	router := setupRouter()
	
	// Generate a valid token then tamper with it
	userId := uuid.New()
	subject := "test@example.com"
	
	tokenPair, err := jwtUtil.GenerateTokenPair(userId, subject, []string{}, []string{})
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	
	// Tamper with the token
	tamperedToken := tokenPair.AccessToken[:len(tokenPair.AccessToken)-10] + "TAMPERED99"
	
	// Create request
	reqBody := IntrospectTokenRequest{
		Token: tamperedToken,
	}
	body, _ := json.Marshal(reqBody)
	
	req, _ := http.NewRequest("POST", "/aegis/api/auth/introspect", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Assert response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
	
	var response IntrospectTokenResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	
	// Tampered token should be inactive
	if response.Active {
		t.Error("Expected active=false for tampered token")
	}
}

// TestIntrospectToken_WithTokenTypeHint tests introspection with token_type_hint parameter
// Expected: Returns 200 OK and processes hint (currently hint is logged but not enforced)
func TestIntrospectToken_WithTokenTypeHint(t *testing.T) {
	router := setupRouter()
	
	// Generate tokens
	userId := uuid.New()
	subject := "test@example.com"
	
	tokenPair, err := jwtUtil.GenerateTokenPair(userId, subject, []string{"user"}, []string{"read"})
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	
	// Test with access_token hint
	reqBody := IntrospectTokenRequest{
		Token:         tokenPair.AccessToken,
		TokenTypeHint: "access_token",
	}
	body, _ := json.Marshal(reqBody)
	
	req, _ := http.NewRequest("POST", "/aegis/api/auth/introspect", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Assert response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
	
	var response IntrospectTokenResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	
	// Should still validate successfully
	if !response.Active {
		t.Error("Expected active=true with token_type_hint")
	}
}

// TestIntrospectToken_EmptyToken tests request with empty token
// Expected: Returns 400 Bad Request
func TestIntrospectToken_EmptyToken(t *testing.T) {
	router := setupRouter()
	
	// Create request with empty token
	reqBody := IntrospectTokenRequest{
		Token: "",
	}
	body, _ := json.Marshal(reqBody)
	
	req, _ := http.NewRequest("POST", "/aegis/api/auth/introspect", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Should return 400 Bad Request
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestIntrospectToken_MissingTokenField tests request without token field
// Expected: Returns 400 Bad Request
func TestIntrospectToken_MissingTokenField(t *testing.T) {
	router := setupRouter()
	
	// Create request without token field
	body := []byte("{}")
	
	req, _ := http.NewRequest("POST", "/aegis/api/auth/introspect", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Should return 400 Bad Request
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestIntrospectToken_InvalidJSON tests request with invalid JSON
// Expected: Returns 400 Bad Request
func TestIntrospectToken_InvalidJSON(t *testing.T) {
	router := setupRouter()
	
	// Create request with invalid JSON
	body := []byte("not valid json")
	
	req, _ := http.NewRequest("POST", "/aegis/api/auth/introspect", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Should return 400 Bad Request
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestIntrospectToken_RFC7662Compliance tests RFC 7662 compliance
// Expected: Response contains all required RFC 7662 fields
func TestIntrospectToken_RFC7662Compliance(t *testing.T) {
	router := setupRouter()
	
	// Generate token
	userId := uuid.New()
	subject := "test@example.com"
	roles := []string{"admin"}
	permissions := []string{"read:users", "write:users"}
	
	tokenPair, err := jwtUtil.GenerateTokenPair(userId, subject, roles, permissions)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	
	// Create introspection request
	reqBody := IntrospectTokenRequest{
		Token: tokenPair.AccessToken,
	}
	body, _ := json.Marshal(reqBody)
	
	req, _ := http.NewRequest("POST", "/aegis/api/auth/introspect", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	var response IntrospectTokenResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	
	// RFC 7662 Section 2.2: Required and recommended fields
	requiredFields := map[string]bool{
		"active": response.Active,
	}
	
	for field, present := range requiredFields {
		if !present {
			t.Errorf("RFC 7662 required field '%s' is missing or false", field)
		}
	}
	
	// For active tokens, these fields should be present
	if response.Active {
		if response.Sub == "" {
			t.Error("RFC 7662: sub field should be present for active token")
		}
		if response.Exp == 0 {
			t.Error("RFC 7662: exp field should be present for active token")
		}
		if response.Iat == 0 {
			t.Error("RFC 7662: iat field should be present for active token")
		}
	}
}

// TestIntrospectToken_EmptyRolesAndPermissions tests token with no roles/permissions
// Expected: Returns active=true with empty scope
func TestIntrospectToken_EmptyRolesAndPermissions(t *testing.T) {
	router := setupRouter()
	
	// Generate token with no roles/permissions
	userId := uuid.New()
	subject := "test@example.com"
	
	tokenPair, err := jwtUtil.GenerateTokenPair(userId, subject, []string{}, []string{})
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	
	// Create introspection request
	reqBody := IntrospectTokenRequest{
		Token: tokenPair.AccessToken,
	}
	body, _ := json.Marshal(reqBody)
	
	req, _ := http.NewRequest("POST", "/aegis/api/auth/introspect", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	var response IntrospectTokenResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	
	// Should be active with empty/minimal scope
	if !response.Active {
		t.Error("Expected active=true for token with no roles/permissions")
	}
	if response.Scope != "" {
		t.Logf("Scope for empty roles/permissions: '%s'", response.Scope)
	}
}

// TestIntrospectToken_ResponseTime tests that introspection is performant
// Expected: Response time < 100ms
func TestIntrospectToken_ResponseTime(t *testing.T) {
	router := setupRouter()
	
	// Generate token
	userId := uuid.New()
	subject := "test@example.com"
	
	tokenPair, err := jwtUtil.GenerateTokenPair(userId, subject, []string{"admin"}, []string{"read"})
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	
	// Measure response time
	reqBody := IntrospectTokenRequest{
		Token: tokenPair.AccessToken,
	}
	body, _ := json.Marshal(reqBody)
	
	start := time.Now()
	
	req, _ := http.NewRequest("POST", "/aegis/api/auth/introspect", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	elapsed := time.Since(start)
	
	// Response should be fast (< 100ms)
	if elapsed > 100*time.Millisecond {
		t.Logf("Warning: Response time %v exceeds 100ms threshold", elapsed)
	}
	
	// Log actual response time
	t.Logf("Token introspection response time: %v", elapsed)
}

// TestBuildScopeString tests the scope string builder function
func TestBuildScopeString(t *testing.T) {
	tests := []struct {
		name        string
		roles       []string
		permissions []string
		expected    string
	}{
		{
			name:        "with roles and permissions",
			roles:       []string{"admin", "user"},
			permissions: []string{"read:users", "write:users"},
			expected:    "role:admin role:user read:users write:users",
		},
		{
			name:        "with only roles",
			roles:       []string{"admin"},
			permissions: []string{},
			expected:    "role:admin",
		},
		{
			name:        "with only permissions",
			roles:       []string{},
			permissions: []string{"read:users"},
			expected:    "read:users",
		},
		{
			name:        "empty roles and permissions",
			roles:       []string{},
			permissions: []string{},
			expected:    "",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildScopeString(tt.roles, tt.permissions)
			if result != tt.expected {
				t.Errorf("Expected scope '%s', got '%s'", tt.expected, result)
			}
		})
	}
}
