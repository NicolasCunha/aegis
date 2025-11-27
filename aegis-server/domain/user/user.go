// Package user provides domain models and business logic for user management,
// including authentication, roles, and permissions.
package user

import (
	"time"
	"github.com/google/uuid"
	"nfcunha/aegis/util/hash"
)

// UserRole represents a role that can be assigned to a user for authorization purposes.
type UserRole string

// Permission represents a specific permission that can be granted to a user.
type Permission string

// User represents a user entity in the system with authentication credentials,
// audit information, roles, and permissions.
type User struct {
	Id	   			uuid.UUID
	Subject			string
	PasswordHash 	string
	Salt			string
	Pepper			string
	CreatedAt		time.Time
	CreatedBy		string
	UpdatedAt		time.Time
	UpdatedBy		string
	AdditionalInfo  map[string]interface{}
	Roles			[]UserRole
	Permissions		[]Permission
}

// CreateUser creates a new User instance with a hashed password.
// A unique ID is generated and the password is securely hashed with a random salt and pepper.
//
// Parameters:
//   - subject: User's subject identifier (typically email or username)
//   - password: Plain text password to be hashed
//   - createdBy: Identifier of who created this user
//
// Returns:
//   - Pointer to the newly created User
func CreateUser(subject string, 
		password string, 
		createdBy string) *User {
	hashOutput := hash.Hash(password)

	return &User{
		Id:             uuid.New(),
		Subject:        subject,
		PasswordHash:   hashOutput.Hash,
		Salt:           hashOutput.Salt,
		Pepper:         hashOutput.Pepper,
		CreatedAt:      time.Now(),
		CreatedBy:      createdBy,
		UpdatedAt:      time.Now(),
		UpdatedBy:      createdBy,
	}
}

// PasswordMatch verifies if the provided password matches the user's stored password hash.
// Uses the stored salt and pepper to recreate the hash for comparison.
//
// Parameters:
//   - password: Plain text password to verify
//
// Returns:
//   - true if the password matches, false otherwise
func (u *User) PasswordMatch(password string) bool {
	return hash.Compare(password, u.Salt, u.Pepper, u.PasswordHash)
}

// UpdatePassword changes the user's password by generating a new hash with fresh salt and pepper.
// Updates the audit fields with the current timestamp and updater identifier.
//
// Parameters:
//   - newPassword: The new plain text password
//   - updatedBy: Identifier of who is updating the password
func (u *User) UpdatePassword(newPassword string, updatedBy string) {
	hashOutput := hash.Hash(newPassword)
	newPasswordHash := hashOutput.Hash
	newSalt := hashOutput.Salt
	newPepper := hashOutput.Pepper
	u.PasswordHash = newPasswordHash
	u.Salt = newSalt
	u.Pepper = newPepper
	u.UpdatedAt = time.Now()
	u.UpdatedBy = updatedBy
}

// UpdateAdditionalInfo updates the user's additional information map and audit fields.
//
// Parameters:
//   - additionalInfo: Map of additional key-value pairs to store
//   - updatedBy: Identifier of who is updating the information
func (u *User) UpdateAdditionalInfo(additionalInfo map[string]interface{}, updatedBy string) {
	u.AdditionalInfo = additionalInfo
	u.UpdatedAt = time.Now()
	u.UpdatedBy = updatedBy
}

// AddRole adds a role to the user if not already present.
// Updates audit fields when a new role is added.
//
// Parameters:
//   - role: The role to add
//   - updatedBy: Identifier of who is adding the role
func (u *User) AddRole(role UserRole, updatedBy string) {
	for _, r := range u.Roles {
		if r == role {
			return
		}
	}
	u.Roles = append(u.Roles, role)
	u.UpdatedAt = time.Now()
	u.UpdatedBy = updatedBy
}

// RemoveRole removes a role from the user if present.
// Updates audit fields when a role is removed.
//
// Parameters:
//   - role: The role to remove
//   - updatedBy: Identifier of who is removing the role
func (u *User) RemoveRole(role UserRole, updatedBy string) {
	for i, r := range u.Roles {
		if r == role {
			u.Roles = append(u.Roles[:i], u.Roles[i+1:]...)
			u.UpdatedAt = time.Now()
			u.UpdatedBy = updatedBy
			return
		}
	}
}

// HasRole checks if the user has a specific role.
//
// Parameters:
//   - role: The role to check for
//
// Returns:
//   - true if the user has the role, false otherwise
func (u *User) HasRole(role UserRole) bool {
	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// AddPermission adds a permission to the user if not already present.
// Updates audit fields when a new permission is added.
//
// Parameters:
//   - permission: The permission to add
//   - updatedBy: Identifier of who is adding the permission
func (u *User) AddPermission(permission Permission, updatedBy string) {
	for _, p := range u.Permissions {
		if p == permission {
			return
		}
	}
	u.Permissions = append(u.Permissions, permission)
	u.UpdatedAt = time.Now()
	u.UpdatedBy = updatedBy
}

// RemovePermission removes a permission from the user if present.
// Updates audit fields when a permission is removed.
//
// Parameters:
//   - permission: The permission to remove
//   - updatedBy: Identifier of who is removing the permission
func (u *User) RemovePermission(permission Permission, updatedBy string) {
	for i, p := range u.Permissions {
		if p == permission {
			u.Permissions = append(u.Permissions[:i], u.Permissions[i+1:]...)
			u.UpdatedAt = time.Now()
			u.UpdatedBy = updatedBy
			return
		}
	}
}

// HasPermission checks if the user has a specific permission.
//
// Parameters:
//   - permission: The permission to check for
//
// Returns:
//   - true if the user has the permission, false otherwise
func (u *User) HasPermission(permission Permission) bool {
	for _, p := range u.Permissions {
		if p == permission {
			return true
		}
	}
	return false
}

