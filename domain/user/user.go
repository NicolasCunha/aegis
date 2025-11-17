package user

import (
	"fmt"
	"time"
	"github.com/google/uuid"
	"nfcunha/aegis/util/hash"
)

type UserRole string
type Permission string

/* User structure */
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

// Create a new user
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

func (u *User) PasswordMatch(password string) bool {
	return hash.Compare(password, u.Salt, u.Pepper, u.PasswordHash)
}

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

func (u *User) UpdateAdditionalInfo(additionalInfo map[string]interface{}, updatedBy string) {
	u.AdditionalInfo = additionalInfo
	u.UpdatedAt = time.Now()
	u.UpdatedBy = updatedBy
}

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

func (u *User) HasRole(role UserRole) bool {
	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}
	return false
}

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

func (u *User) HasPermission(permission Permission) bool {
	for _, p := range u.Permissions {
		if p == permission {
			return true
		}
	}
	return false
}

func (u *User) FormatWithAllFields() string {
	return fmt.Sprintf("User{Id: %s, \nSubject: %s, \nPasswordHash: %s, \nSalt: %s, \nPepper: %s, \nCreatedAt: %s, \nCreatedBy: %s, \nUpdatedAt: %s, \nUpdatedBy: %s, \nAdditionalInfo: %v, \nRoles: %v, \nPermissions: %v}",
		u.Id, u.Subject, u.PasswordHash, u.Salt, u.Pepper,
		u.CreatedAt, u.CreatedBy, u.UpdatedAt, u.UpdatedBy,
		u.AdditionalInfo, u.Roles, u.Permissions)
}