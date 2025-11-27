// Package user provides HTTP REST API endpoints for user management operations.
// Includes registration, authentication, CRUD operations, and password management.
package user

import (
	"log"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	userService "nfcunha/aegis/domain/user"
	"nfcunha/aegis/util/jwt"
)

type RegisterRequest struct {
	Subject     string   `json:"subject" binding:"required"`
	Password    string   `json:"password" binding:"required,min=8"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}

type LoginRequest struct {
	Subject  string `json:"subject" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserRequest struct {
	Subject     string   `json:"subject"`
	Password    string   `json:"password,omitempty"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type AddRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

type AddPermissionRequest struct {
	Permission string `json:"permission" binding:"required"`
}

type UserResponse struct {
	Id          string                     `json:"id"`
	Subject     string                     `json:"subject"`
	CreatedAt   time.Time                  `json:"created_at"`
	CreatedBy   string                     `json:"created_by"`
	UpdatedAt   time.Time                  `json:"updated_at"`
	UpdatedBy   string                     `json:"updated_by"`
	Roles       []userService.UserRole     `json:"roles"`
	Permissions []userService.Permission   `json:"permissions"`
}

type LoginResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresAt    time.Time    `json:"expires_at"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at"`
}

// RegisterApi registers all user-related HTTP routes with the Gin router.
// Endpoints include register, login, list, get, update, delete, and change password.
//
// Parameters:
//   - router: The Gin engine to register routes with
func RegisterApi(router *gin.Engine) {
	users := router.Group("/users")
	{
		users.POST("/register", registerUser)
		users.POST("/login", loginUser)
		users.POST("/refresh", refreshToken)
		users.GET("", listUsers)
		users.GET("/:id", getUser)
		users.PUT("/:id", updateUser)
		users.DELETE("/:id", deleteUser)
		users.POST("/:id/password", changePassword)
		users.POST("/:id/roles", addRoleToUser)
		users.DELETE("/:id/roles/:role", removeRoleFromUser)
		users.POST("/:id/permissions", addPermissionToUser)
		users.DELETE("/:id/permissions/:permission", removePermissionFromUser)
	}
}

func registerUser(c *gin.Context) {
	log.Println("POST /users/register - Register user request received")
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists
	if userService.ExistsUserBySubject(req.Subject) {
		log.Printf("User already exists: %s", req.Subject)
		c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
		return
	}

	// Create user
	user := userService.CreateUser(req.Subject, req.Password, "system")
	
	// Add roles
	for _, role := range req.Roles {
		user.Roles = append(user.Roles, userService.UserRole(role))
	}
	
	// Add permissions
	for _, permission := range req.Permissions {
		user.Permissions = append(user.Permissions, userService.Permission(permission))
	}

	// Persist user
	userService.PersistUser(user)

	log.Printf("User registered successfully: %s", user.Subject)
	c.JSON(http.StatusCreated, toUserResponse(user))
}

func loginUser(c *gin.Context) {
	log.Println("POST /users/login - Login request received")
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user by subject
	user := userService.GetUserBySubject(req.Subject)
	if user == nil {
		log.Printf("Login failed: user not found - %s", req.Subject)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Check password
	if !user.PasswordMatch(req.Password) {
		log.Printf("Login failed: invalid password - %s", req.Subject)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Generate token
	roles := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = string(role)
	}

	permissions := make([]string, len(user.Permissions))
	for i, permission := range user.Permissions {
		permissions[i] = string(permission)
	}

	tokenPair, err := jwt.GenerateTokenPair(user.Id, user.Subject, roles, permissions)
	if err != nil {
		log.Printf("Failed to generate tokens for user %s: %v", req.Subject, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate tokens"})
		return
	}

	log.Printf("User logged in successfully: %s", user.Subject)
	c.JSON(http.StatusOK, LoginResponse{
		User:             toUserResponse(user),
		AccessToken:      tokenPair.AccessToken,
		RefreshToken:     tokenPair.RefreshToken,
		ExpiresAt:        tokenPair.ExpiresAt,
		RefreshExpiresAt: tokenPair.RefreshExpiresAt,
	})
}

func refreshToken(c *gin.Context) {
	log.Println("POST /users/refresh - Refresh token request received")
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate refresh token
	claims, err := jwt.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		log.Printf("Invalid refresh token: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	// Get user from claims
	userId, err := uuid.Parse(claims.UserId)
	if err != nil {
		log.Printf("Invalid user ID in token: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	user := userService.GetUserById(userId)
	if user == nil {
		log.Printf("User not found for token refresh: %s", claims.UserId)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	// Generate new token pair
	roles := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = string(role)
	}

	permissions := make([]string, len(user.Permissions))
	for i, permission := range user.Permissions {
		permissions[i] = string(permission)
	}

	tokenPair, err := jwt.GenerateTokenPair(user.Id, user.Subject, roles, permissions)
	if err != nil {
		log.Printf("Failed to generate new tokens for user %s: %v", user.Subject, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate tokens"})
		return
	}

	log.Printf("Token refreshed successfully for user: %s", user.Subject)
	c.JSON(http.StatusOK, LoginResponse{
		User:             toUserResponse(user),
		AccessToken:      tokenPair.AccessToken,
		RefreshToken:     tokenPair.RefreshToken,
		ExpiresAt:        tokenPair.ExpiresAt,
		RefreshExpiresAt: tokenPair.RefreshExpiresAt,
	})
}

func listUsers(c *gin.Context) {
	log.Println("GET /users - List users request received")
	users := userService.ListUsers()
	response := make([]UserResponse, len(users))
	for i, user := range users {
		response[i] = toUserResponse(user)
	}
	log.Printf("Returning %d users", len(response))
	c.JSON(http.StatusOK, response)
}

func getUser(c *gin.Context) {
	idStr := c.Param("id")
	log.Printf("GET /users/%s - Get user request received", idStr)
	userId, err := uuid.Parse(idStr)
	if err != nil {
		log.Printf("Invalid user ID: %s", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	user := userService.GetUserById(userId)
	if user == nil {
		log.Printf("User not found: %s", idStr)
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	log.Printf("Returning user: %s", user.Subject)
	c.JSON(http.StatusOK, toUserResponse(user))
}

func updateUser(c *gin.Context) {
	idStr := c.Param("id")
	userId, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	user := userService.GetUserById(userId)
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update subject if provided
	if req.Subject != "" && req.Subject != user.Subject {
		// Check if new subject already exists
		if userService.ExistsUserBySubject(req.Subject) {
			c.JSON(http.StatusConflict, gin.H{"error": "subject already exists"})
			return
		}
		user.Subject = req.Subject
	}

	// Update password if provided
	if req.Password != "" {
		user.UpdatePassword(req.Password, "system")
	}

	// Update roles
	user.Roles = make([]userService.UserRole, len(req.Roles))
	for i, role := range req.Roles {
		user.Roles[i] = userService.UserRole(role)
	}

	// Update permissions
	user.Permissions = make([]userService.Permission, len(req.Permissions))
	for i, permission := range req.Permissions {
		user.Permissions[i] = userService.Permission(permission)
	}

	user.UpdatedAt = time.Now()
	user.UpdatedBy = "system"

	userService.PersistUser(user)

	c.JSON(http.StatusOK, toUserResponse(user))
}

func deleteUser(c *gin.Context) {
	idStr := c.Param("id")
	log.Printf("DELETE /users/%s - Delete user request received", idStr)
	userId, err := uuid.Parse(idStr)
	if err != nil {
		log.Printf("Invalid user ID: %s", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	user := userService.GetUserById(userId)
	if user == nil {
		log.Printf("User not found: %s", idStr)
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	userService.DeleteUser(userId)

	log.Printf("User deleted: %s", user.Subject)
	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}

func changePassword(c *gin.Context) {
	idStr := c.Param("id")
	userId, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	user := userService.GetUserById(userId)
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify old password
	if !user.PasswordMatch(req.OldPassword) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid old password"})
		return
	}

	// Update password
	user.UpdatePassword(req.NewPassword, "system")
	userService.UpdateUser(user)

	c.JSON(http.StatusOK, gin.H{"message": "password changed successfully"})
}

func addRoleToUser(c *gin.Context) {
	idStr := c.Param("id")
	log.Printf("POST /users/%s/roles - Add role to user request received", idStr)
	userId, err := uuid.Parse(idStr)
	if err != nil {
		log.Printf("Invalid user ID: %s", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	user := userService.GetUserById(userId)
	if user == nil {
		log.Printf("User not found: %s", idStr)
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	var req AddRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already has the role
	role := userService.UserRole(req.Role)
	if user.HasRole(role) {
		log.Printf("User %s already has role: %s", user.Subject, req.Role)
		c.JSON(http.StatusConflict, gin.H{"error": "user already has this role"})
		return
	}

	// Add role
	user.AddRole(role, "system")
	userService.AddUserRole(user, role)

	log.Printf("Role %s added to user %s", req.Role, user.Subject)
	c.JSON(http.StatusOK, toUserResponse(user))
}

func removeRoleFromUser(c *gin.Context) {
	idStr := c.Param("id")
	roleStr := c.Param("role")
	log.Printf("DELETE /users/%s/roles/%s - Remove role from user request received", idStr, roleStr)
	userId, err := uuid.Parse(idStr)
	if err != nil {
		log.Printf("Invalid user ID: %s", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	user := userService.GetUserById(userId)
	if user == nil {
		log.Printf("User not found: %s", idStr)
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Check if user has the role
	role := userService.UserRole(roleStr)
	if !user.HasRole(role) {
		log.Printf("User %s does not have role: %s", user.Subject, roleStr)
		c.JSON(http.StatusNotFound, gin.H{"error": "user does not have this role"})
		return
	}

	// Remove role
	user.RemoveRole(role, "system")
	userService.RemoveUserRole(user, role)

	log.Printf("Role %s removed from user %s", roleStr, user.Subject)
	c.JSON(http.StatusOK, toUserResponse(user))
}

func addPermissionToUser(c *gin.Context) {
	idStr := c.Param("id")
	log.Printf("POST /users/%s/permissions - Add permission to user request received", idStr)
	userId, err := uuid.Parse(idStr)
	if err != nil {
		log.Printf("Invalid user ID: %s", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	user := userService.GetUserById(userId)
	if user == nil {
		log.Printf("User not found: %s", idStr)
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	var req AddPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already has the permission
	permission := userService.Permission(req.Permission)
	if user.HasPermission(permission) {
		log.Printf("User %s already has permission: %s", user.Subject, req.Permission)
		c.JSON(http.StatusConflict, gin.H{"error": "user already has this permission"})
		return
	}

	// Add permission
	user.AddPermission(permission, "system")
	userService.AddUserPermission(user, permission)

	log.Printf("Permission %s added to user %s", req.Permission, user.Subject)
	c.JSON(http.StatusOK, toUserResponse(user))
}

func removePermissionFromUser(c *gin.Context) {
	idStr := c.Param("id")
	permissionStr := c.Param("permission")
	log.Printf("DELETE /users/%s/permissions/%s - Remove permission from user request received", idStr, permissionStr)
	userId, err := uuid.Parse(idStr)
	if err != nil {
		log.Printf("Invalid user ID: %s", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	user := userService.GetUserById(userId)
	if user == nil {
		log.Printf("User not found: %s", idStr)
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Check if user has the permission
	permission := userService.Permission(permissionStr)
	if !user.HasPermission(permission) {
		log.Printf("User %s does not have permission: %s", user.Subject, permissionStr)
		c.JSON(http.StatusNotFound, gin.H{"error": "user does not have this permission"})
		return
	}

	// Remove permission
	user.RemovePermission(permission, "system")
	userService.RemoveUserPermission(user, permission)

	log.Printf("Permission %s removed from user %s", permissionStr, user.Subject)
	c.JSON(http.StatusOK, toUserResponse(user))
}

// toUserResponse converts a domain User model to an API UserResponse.
// Excludes sensitive fields like password hashes.
//
// Parameters:
//   - user: The domain user to convert
//
// Returns:
//   - UserResponse containing safe user data for API responses
func toUserResponse(user *userService.User) UserResponse {
	return UserResponse{
		Id:          user.Id.String(),
		Subject:     user.Subject,
		CreatedAt:   user.CreatedAt,
		CreatedBy:   user.CreatedBy,
		UpdatedAt:   user.UpdatedAt,
		UpdatedBy:   user.UpdatedBy,
		Roles:       user.Roles,
		Permissions: user.Permissions,
	}
}