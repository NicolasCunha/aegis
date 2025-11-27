package user

import (
	"testing"
	"time"
)

// TestCreateUser tests user creation with password hashing
func TestCreateUser(t *testing.T) {
	subject := "test@example.com"
	password := "password123"
	createdBy := "admin"
	
	user := CreateUser(subject, password, createdBy)
	
	if user == nil {
		t.Fatal("CreateUser should not return nil")
	}
	if user.Id.String() == "00000000-0000-0000-0000-000000000000" {
		t.Error("User ID should not be nil UUID")
	}
	if user.Subject != subject {
		t.Errorf("Expected subject %s, got %s", subject, user.Subject)
	}
	if user.PasswordHash == "" {
		t.Error("Password hash should not be empty")
	}
	if user.Salt == "" {
		t.Error("Salt should not be empty")
	}
	if user.Pepper == "" {
		t.Error("Pepper should not be empty")
	}
	if user.CreatedBy != createdBy {
		t.Errorf("Expected createdBy %s, got %s", createdBy, user.CreatedBy)
	}
	if user.UpdatedBy != createdBy {
		t.Errorf("Expected updatedBy %s, got %s", createdBy, user.UpdatedBy)
	}
	if user.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
	if user.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should not be zero")
	}
}

// TestCreateUser_UniqueIds tests that each user gets a unique ID
func TestCreateUser_UniqueIds(t *testing.T) {
	user1 := CreateUser("user1@example.com", "password", "admin")
	user2 := CreateUser("user2@example.com", "password", "admin")
	
	if user1.Id == user2.Id {
		t.Error("Each user should have a unique ID")
	}
}

// TestCreateUser_UniqueHashComponents tests unique salt and pepper per user
func TestCreateUser_UniqueHashComponents(t *testing.T) {
	password := "samepassword"
	user1 := CreateUser("user1@example.com", password, "admin")
	user2 := CreateUser("user2@example.com", password, "admin")
	
	// Same password should produce different hashes due to unique salt/pepper
	if user1.PasswordHash == user2.PasswordHash {
		t.Error("Same password should produce different hashes for different users")
	}
	if user1.Salt == user2.Salt {
		t.Error("Each user should have unique salt")
	}
	if user1.Pepper == user2.Pepper {
		t.Error("Each user should have unique pepper")
	}
}

// TestPasswordMatch_ValidPassword tests successful password verification
func TestPasswordMatch_ValidPassword(t *testing.T) {
	password := "password123"
	user := CreateUser("test@example.com", password, "admin")
	
	if !user.PasswordMatch(password) {
		t.Error("PasswordMatch should return true for correct password")
	}
}

// TestPasswordMatch_InvalidPassword tests failed password verification
func TestPasswordMatch_InvalidPassword(t *testing.T) {
	password := "password123"
	user := CreateUser("test@example.com", password, "admin")
	
	if user.PasswordMatch("wrongpassword") {
		t.Error("PasswordMatch should return false for incorrect password")
	}
}

// TestPasswordMatch_EmptyPassword tests password verification with empty password
func TestPasswordMatch_EmptyPassword(t *testing.T) {
	user := CreateUser("test@example.com", "password123", "admin")
	
	if user.PasswordMatch("") {
		t.Error("PasswordMatch should return false for empty password")
	}
}

// TestPasswordMatch_CaseSensitive tests that password matching is case-sensitive
func TestPasswordMatch_CaseSensitive(t *testing.T) {
	password := "Password123"
	user := CreateUser("test@example.com", password, "admin")
	
	if user.PasswordMatch("password123") {
		t.Error("PasswordMatch should be case-sensitive")
	}
}

// TestUpdatePassword tests password update functionality
func TestUpdatePassword(t *testing.T) {
	user := CreateUser("test@example.com", "oldpassword", "admin")
	oldHash := user.PasswordHash
	oldSalt := user.Salt
	oldPepper := user.Pepper
	oldUpdatedAt := user.UpdatedAt
	
	time.Sleep(1 * time.Millisecond) // Ensure timestamp difference
	newPassword := "newpassword"
	updatedBy := "user"
	user.UpdatePassword(newPassword, updatedBy)
	
	// Verify password was changed
	if user.PasswordHash == oldHash {
		t.Error("Password hash should change after update")
	}
	if user.Salt == oldSalt {
		t.Error("Salt should change after password update")
	}
	if user.Pepper == oldPepper {
		t.Error("Pepper should change after password update")
	}
	
	// Verify new password works
	if !user.PasswordMatch(newPassword) {
		t.Error("New password should match after update")
	}
	
	// Verify old password no longer works
	if user.PasswordMatch("oldpassword") {
		t.Error("Old password should not match after update")
	}
	
	// Verify audit fields
	if user.UpdatedBy != updatedBy {
		t.Errorf("Expected updatedBy %s, got %s", updatedBy, user.UpdatedBy)
	}
	if !user.UpdatedAt.After(oldUpdatedAt) {
		t.Error("UpdatedAt should be updated")
	}
}

// TestUpdateAdditionalInfo tests updating additional information
func TestUpdateAdditionalInfo(t *testing.T) {
	user := CreateUser("test@example.com", "password", "admin")
	oldUpdatedAt := user.UpdatedAt
	
	time.Sleep(1 * time.Millisecond)
	additionalInfo := map[string]interface{}{
		"department": "Engineering",
		"level":      5,
	}
	updatedBy := "admin"
	user.UpdateAdditionalInfo(additionalInfo, updatedBy)
	
	if user.AdditionalInfo == nil {
		t.Error("AdditionalInfo should not be nil")
	}
	if user.AdditionalInfo["department"] != "Engineering" {
		t.Error("AdditionalInfo should contain department")
	}
	if user.AdditionalInfo["level"] != 5 {
		t.Error("AdditionalInfo should contain level")
	}
	if user.UpdatedBy != updatedBy {
		t.Errorf("Expected updatedBy %s, got %s", updatedBy, user.UpdatedBy)
	}
	if !user.UpdatedAt.After(oldUpdatedAt) {
		t.Error("UpdatedAt should be updated")
	}
}

// TestAddRole tests adding a role to a user
func TestAddRole(t *testing.T) {
	user := CreateUser("test@example.com", "password", "admin")
	oldUpdatedAt := user.UpdatedAt
	
	time.Sleep(1 * time.Millisecond)
	role := UserRole("admin")
	updatedBy := "system"
	user.AddRole(role, updatedBy)
	
	if len(user.Roles) != 1 {
		t.Errorf("Expected 1 role, got %d", len(user.Roles))
	}
	if user.Roles[0] != role {
		t.Errorf("Expected role %s, got %s", role, user.Roles[0])
	}
	if user.UpdatedBy != updatedBy {
		t.Errorf("Expected updatedBy %s, got %s", updatedBy, user.UpdatedBy)
	}
	if !user.UpdatedAt.After(oldUpdatedAt) {
		t.Error("UpdatedAt should be updated when adding role")
	}
}

// TestAddRole_Duplicate tests that duplicate roles are not added
func TestAddRole_Duplicate(t *testing.T) {
	user := CreateUser("test@example.com", "password", "admin")
	role := UserRole("admin")
	
	user.AddRole(role, "system")
	user.AddRole(role, "system") // Try to add again
	
	if len(user.Roles) != 1 {
		t.Errorf("Expected 1 role, got %d", len(user.Roles))
	}
}

// TestAddRole_Multiple tests adding multiple roles
func TestAddRole_Multiple(t *testing.T) {
	user := CreateUser("test@example.com", "password", "admin")
	
	user.AddRole(UserRole("admin"), "system")
	user.AddRole(UserRole("moderator"), "system")
	user.AddRole(UserRole("user"), "system")
	
	if len(user.Roles) != 3 {
		t.Errorf("Expected 3 roles, got %d", len(user.Roles))
	}
}

// TestHasRole_True tests role check returns true when user has role
func TestHasRole_True(t *testing.T) {
	user := CreateUser("test@example.com", "password", "admin")
	role := UserRole("admin")
	user.AddRole(role, "system")
	
	if !user.HasRole(role) {
		t.Error("HasRole should return true when user has the role")
	}
}

// TestHasRole_False tests role check returns false when user doesn't have role
func TestHasRole_False(t *testing.T) {
	user := CreateUser("test@example.com", "password", "admin")
	
	if user.HasRole(UserRole("admin")) {
		t.Error("HasRole should return false when user doesn't have the role")
	}
}

// TestRemoveRole tests removing a role from a user
func TestRemoveRole(t *testing.T) {
	user := CreateUser("test@example.com", "password", "admin")
	role := UserRole("admin")
	user.AddRole(role, "system")
	oldUpdatedAt := user.UpdatedAt
	
	time.Sleep(1 * time.Millisecond)
	updatedBy := "system"
	user.RemoveRole(role, updatedBy)
	
	if len(user.Roles) != 0 {
		t.Errorf("Expected 0 roles, got %d", len(user.Roles))
	}
	if user.HasRole(role) {
		t.Error("User should not have the removed role")
	}
	if user.UpdatedBy != updatedBy {
		t.Errorf("Expected updatedBy %s, got %s", updatedBy, user.UpdatedBy)
	}
	if !user.UpdatedAt.After(oldUpdatedAt) {
		t.Error("UpdatedAt should be updated when removing role")
	}
}

// TestRemoveRole_NonExistent tests removing a role that user doesn't have
func TestRemoveRole_NonExistent(t *testing.T) {
	user := CreateUser("test@example.com", "password", "admin")
	user.AddRole(UserRole("admin"), "system")
	oldUpdatedAt := user.UpdatedAt
	
	time.Sleep(1 * time.Millisecond)
	user.RemoveRole(UserRole("nonexistent"), "system")
	
	// Should not change anything
	if len(user.Roles) != 1 {
		t.Errorf("Expected 1 role, got %d", len(user.Roles))
	}
	if user.UpdatedAt.After(oldUpdatedAt) {
		t.Error("UpdatedAt should not change when removing non-existent role")
	}
}

// TestRemoveRole_FromMultiple tests removing one role when user has multiple
func TestRemoveRole_FromMultiple(t *testing.T) {
	user := CreateUser("test@example.com", "password", "admin")
	user.AddRole(UserRole("admin"), "system")
	user.AddRole(UserRole("moderator"), "system")
	user.AddRole(UserRole("user"), "system")
	
	user.RemoveRole(UserRole("moderator"), "system")
	
	if len(user.Roles) != 2 {
		t.Errorf("Expected 2 roles, got %d", len(user.Roles))
	}
	if user.HasRole(UserRole("moderator")) {
		t.Error("User should not have removed role")
	}
	if !user.HasRole(UserRole("admin")) || !user.HasRole(UserRole("user")) {
		t.Error("Other roles should remain")
	}
}

// TestAddPermission tests adding a permission to a user
func TestAddPermission(t *testing.T) {
	user := CreateUser("test@example.com", "password", "admin")
	oldUpdatedAt := user.UpdatedAt
	
	time.Sleep(1 * time.Millisecond)
	permission := Permission("read")
	updatedBy := "system"
	user.AddPermission(permission, updatedBy)
	
	if len(user.Permissions) != 1 {
		t.Errorf("Expected 1 permission, got %d", len(user.Permissions))
	}
	if user.Permissions[0] != permission {
		t.Errorf("Expected permission %s, got %s", permission, user.Permissions[0])
	}
	if user.UpdatedBy != updatedBy {
		t.Errorf("Expected updatedBy %s, got %s", updatedBy, user.UpdatedBy)
	}
	if !user.UpdatedAt.After(oldUpdatedAt) {
		t.Error("UpdatedAt should be updated when adding permission")
	}
}

// TestAddPermission_Duplicate tests that duplicate permissions are not added
func TestAddPermission_Duplicate(t *testing.T) {
	user := CreateUser("test@example.com", "password", "admin")
	permission := Permission("read")
	
	user.AddPermission(permission, "system")
	user.AddPermission(permission, "system") // Try to add again
	
	if len(user.Permissions) != 1 {
		t.Errorf("Expected 1 permission, got %d", len(user.Permissions))
	}
}

// TestAddPermission_Multiple tests adding multiple permissions
func TestAddPermission_Multiple(t *testing.T) {
	user := CreateUser("test@example.com", "password", "admin")
	
	user.AddPermission(Permission("read"), "system")
	user.AddPermission(Permission("write"), "system")
	user.AddPermission(Permission("delete"), "system")
	
	if len(user.Permissions) != 3 {
		t.Errorf("Expected 3 permissions, got %d", len(user.Permissions))
	}
}

// TestHasPermission_True tests permission check returns true when user has permission
func TestHasPermission_True(t *testing.T) {
	user := CreateUser("test@example.com", "password", "admin")
	permission := Permission("read")
	user.AddPermission(permission, "system")
	
	if !user.HasPermission(permission) {
		t.Error("HasPermission should return true when user has the permission")
	}
}

// TestHasPermission_False tests permission check returns false when user doesn't have permission
func TestHasPermission_False(t *testing.T) {
	user := CreateUser("test@example.com", "password", "admin")
	
	if user.HasPermission(Permission("read")) {
		t.Error("HasPermission should return false when user doesn't have the permission")
	}
}

// TestRemovePermission tests removing a permission from a user
func TestRemovePermission(t *testing.T) {
	user := CreateUser("test@example.com", "password", "admin")
	permission := Permission("read")
	user.AddPermission(permission, "system")
	oldUpdatedAt := user.UpdatedAt
	
	time.Sleep(1 * time.Millisecond)
	updatedBy := "system"
	user.RemovePermission(permission, updatedBy)
	
	if len(user.Permissions) != 0 {
		t.Errorf("Expected 0 permissions, got %d", len(user.Permissions))
	}
	if user.HasPermission(permission) {
		t.Error("User should not have the removed permission")
	}
	if user.UpdatedBy != updatedBy {
		t.Errorf("Expected updatedBy %s, got %s", updatedBy, user.UpdatedBy)
	}
	if !user.UpdatedAt.After(oldUpdatedAt) {
		t.Error("UpdatedAt should be updated when removing permission")
	}
}

// TestRemovePermission_NonExistent tests removing a permission that user doesn't have
func TestRemovePermission_NonExistent(t *testing.T) {
	user := CreateUser("test@example.com", "password", "admin")
	user.AddPermission(Permission("read"), "system")
	oldUpdatedAt := user.UpdatedAt
	
	time.Sleep(1 * time.Millisecond)
	user.RemovePermission(Permission("nonexistent"), "system")
	
	// Should not change anything
	if len(user.Permissions) != 1 {
		t.Errorf("Expected 1 permission, got %d", len(user.Permissions))
	}
	if user.UpdatedAt.After(oldUpdatedAt) {
		t.Error("UpdatedAt should not change when removing non-existent permission")
	}
}

// TestRemovePermission_FromMultiple tests removing one permission when user has multiple
func TestRemovePermission_FromMultiple(t *testing.T) {
	user := CreateUser("test@example.com", "password", "admin")
	user.AddPermission(Permission("read"), "system")
	user.AddPermission(Permission("write"), "system")
	user.AddPermission(Permission("delete"), "system")
	
	user.RemovePermission(Permission("write"), "system")
	
	if len(user.Permissions) != 2 {
		t.Errorf("Expected 2 permissions, got %d", len(user.Permissions))
	}
	if user.HasPermission(Permission("write")) {
		t.Error("User should not have removed permission")
	}
	if !user.HasPermission(Permission("read")) || !user.HasPermission(Permission("delete")) {
		t.Error("Other permissions should remain")
	}
}
