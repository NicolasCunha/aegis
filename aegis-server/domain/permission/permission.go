// Package permission provides domain models and business logic for permission management.
// Permissions represent specific actions or access rights that can be granted to users.
package permission

import (
	"time"
)

// Permission represents a specific access right or action that can be granted to users.
// Permissions are identified by their name and include audit information.
type Permission struct {
	Name        string
	Description string
	CreatedAt   time.Time
	CreatedBy   string
	UpdatedAt   time.Time
	UpdatedBy   string
}

// CreatePermission creates a new Permission instance with the specified name and description.
// Initializes timestamps with the current time.
//
// Parameters:
//   - name: Unique identifier for the permission
//   - description: Human-readable description of what the permission allows
//   - createdBy: Identifier of who created this permission
//
// Returns:
//   - Pointer to the newly created Permission
func CreatePermission(name string, description string, createdBy string) *Permission {
	return &Permission{
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		CreatedBy:   createdBy,
		UpdatedAt:   time.Now(),
		UpdatedBy:   createdBy,
	}
}

// Update modifies the permission's description and updates audit fields.
//
// Parameters:
//   - description: New description for the permission
//   - updatedBy: Identifier of who is updating the permission
func (p *Permission) Update(description string, updatedBy string) {
	p.Description = description
	p.UpdatedAt = time.Now()
	p.UpdatedBy = updatedBy
}


