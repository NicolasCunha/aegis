package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"github.com/gin-gonic/gin"
	"nfcunha/aegis/database"
	userService "nfcunha/aegis/domain/user"
)

func setupTestDB() {
	database.SetTestMode()
	database.Migrate()
}

func teardownTestDB() {
	// Close any open connections
	db, _ := database.OpenConnection()
	if db != nil {
		db.Close()
	}
	
	// Remove test database
	os.Remove("aegis-test.db")
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	RegisterApi(router)
	return router
}

func TestMain(m *testing.M) {
	// Setup
	setupTestDB()
	
	// Run tests
	code := m.Run()
	
	// Teardown
	teardownTestDB()
	
	// Exit
	os.Exit(code)
}

func TestRegisterUser_Success(t *testing.T) {
	router := setupRouter()
	
	reqBody := RegisterRequest{
		Subject:     "register1@example.com",
		Password:    "password123",
		Roles:       []string{"user"},
		Permissions: []string{"read"},
	}
	body, _ := json.Marshal(reqBody)
	
	req, _ := http.NewRequest("POST", "/users/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusCreated, w.Code, w.Body.String())
	}
	
	var response UserResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	
	if response.Subject != reqBody.Subject {
		t.Errorf("Expected subject %s, got %s", reqBody.Subject, response.Subject)
	}
}

func TestRegisterUser_InvalidPassword(t *testing.T) {
	router := setupRouter()
	
	reqBody := RegisterRequest{
		Subject:  "register2@example.com",
		Password: "short",
	}
	body, _ := json.Marshal(reqBody)
	
	req, _ := http.NewRequest("POST", "/users/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestRegisterUser_DuplicateSubject(t *testing.T) {
	router := setupRouter()
	
	// Register first user
	reqBody1 := RegisterRequest{
		Subject:  "register3@example.com",
		Password: "password123",
	}
	body1, _ := json.Marshal(reqBody1)
	req1, _ := http.NewRequest("POST", "/users/register", bytes.NewBuffer(body1))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)
	
	// Try to register duplicate
	reqBody2 := RegisterRequest{
		Subject:  "register3@example.com",
		Password: "password456",
	}
	body2, _ := json.Marshal(reqBody2)
	req2, _ := http.NewRequest("POST", "/users/register", bytes.NewBuffer(body2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	
	if w2.Code != http.StatusConflict {
		t.Errorf("Expected status %d, got %d", http.StatusConflict, w2.Code)
	}
}

func TestLoginUser_Success(t *testing.T) {
	router := setupRouter()
	
	// Register user first
	password := "password123"
	regBody := RegisterRequest{
		Subject:  "login1@example.com",
		Password: password,
	}
	regBodyJSON, _ := json.Marshal(regBody)
	regReq, _ := http.NewRequest("POST", "/users/register", bytes.NewBuffer(regBodyJSON))
	regReq.Header.Set("Content-Type", "application/json")
	regW := httptest.NewRecorder()
	router.ServeHTTP(regW, regReq)
	
	// Login
	loginBody := LoginRequest{
		Subject:  "login1@example.com",
		Password: password,
	}
	loginBodyJSON, _ := json.Marshal(loginBody)
	
	req, _ := http.NewRequest("POST", "/users/login", bytes.NewBuffer(loginBodyJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
	}
	
	var response LoginResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	
	if response.AccessToken == "" {
		t.Error("AccessToken should not be empty")
	}
	if response.RefreshToken == "" {
		t.Error("RefreshToken should not be empty")
	}
}

func TestLoginUser_WrongPassword(t *testing.T) {
	router := setupRouter()
	
	// Register user
	regBody := RegisterRequest{
		Subject:  "login2@example.com",
		Password: "password123",
	}
	regBodyJSON, _ := json.Marshal(regBody)
	regReq, _ := http.NewRequest("POST", "/users/register", bytes.NewBuffer(regBodyJSON))
	regReq.Header.Set("Content-Type", "application/json")
	regW := httptest.NewRecorder()
	router.ServeHTTP(regW, regReq)
	
	// Login with wrong password
	loginBody := LoginRequest{
		Subject:  "login2@example.com",
		Password: "wrongpassword",
	}
	loginBodyJSON, _ := json.Marshal(loginBody)
	
	req, _ := http.NewRequest("POST", "/users/login", bytes.NewBuffer(loginBodyJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestLoginUser_NonExistentUser(t *testing.T) {
	router := setupRouter()
	
	loginBody := LoginRequest{
		Subject:  "nonexistent@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(loginBody)
	
	req, _ := http.NewRequest("POST", "/users/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestListUsers(t *testing.T) {
	router := setupRouter()
	
	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
	
	var response []UserResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	
	// Should have users from previous tests
	if len(response) < 1 {
		t.Errorf("Expected at least 1 user, got %d", len(response))
	}
}

func TestGetUser_Success(t *testing.T) {
	router := setupRouter()
	
	// Create user
	user := userService.CreateUser("getuser1@example.com", "password123", "system")
	userService.SaveUser(user)
	
	req, _ := http.NewRequest("GET", "/users/"+user.Id.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
	
	var response UserResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	
	if response.Id != user.Id.String() {
		t.Errorf("Expected user ID %s, got %s", user.Id.String(), response.Id)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	router := setupRouter()
	
	req, _ := http.NewRequest("GET", "/users/00000000-0000-0000-0000-000000000000", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestDeleteUser_Success(t *testing.T) {
	router := setupRouter()
	
	// Create user
	user := userService.CreateUser("deleteuser1@example.com", "password123", "system")
	userService.SaveUser(user)
	
	req, _ := http.NewRequest("DELETE", "/users/"+user.Id.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
	
	// Verify user is deleted
	deletedUser := userService.GetUserById(user.Id)
	if deletedUser != nil {
		t.Error("User should be deleted")
	}
}

func TestChangePassword_Success(t *testing.T) {
	router := setupRouter()
	
	// Create user
	oldPassword := "oldpassword123"
	user := userService.CreateUser("changepass1@example.com", oldPassword, "system")
	userService.SaveUser(user)
	
	reqBody := ChangePasswordRequest{
		OldPassword: oldPassword,
		NewPassword: "newpassword123",
	}
	body, _ := json.Marshal(reqBody)
	
	req, _ := http.NewRequest("POST", "/users/"+user.Id.String()+"/password", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
	}
	
	// Verify new password works
	updatedUser := userService.GetUserById(user.Id)
	if !updatedUser.PasswordMatch(reqBody.NewPassword) {
		t.Error("New password should match")
	}
	if updatedUser.PasswordMatch(oldPassword) {
		t.Error("Old password should not match")
	}
}

func TestChangePassword_WrongOldPassword(t *testing.T) {
	router := setupRouter()
	
	// Create user
	user := userService.CreateUser("changepass2@example.com", "password123", "system")
	userService.SaveUser(user)
	
	reqBody := ChangePasswordRequest{
		OldPassword: "wrongoldpassword",
		NewPassword: "newpassword123",
	}
	body, _ := json.Marshal(reqBody)
	
	req, _ := http.NewRequest("POST", "/users/"+user.Id.String()+"/password", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}
