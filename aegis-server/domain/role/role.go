// Package role provides domain models and business logic for role management.
// Roles can be assigned to users for authorization purposes.
package role

import (
	"time"
)

// Role represents a role that can be assigned to users for authorization.
// Roles are identified by their name and include audit information.
type Role struct {
	Name        string
	Description string
	CreatedAt   time.Time
	CreatedBy   string
	UpdatedAt   time.Time
	UpdatedBy   string
}

// CreateRole creates a new Role instance with the specified name and description.
// Initializes timestamps with the current time.
//
// Parameters:
//   - name: Unique identifier for the role
//   - description: Human-readable description of the role's purpose
//   - createdBy: Identifier of who created this role
//
// Returns:
//   - Pointer to the newly created Role
func CreateRole(name string, description string, createdBy string) *Role {
	return &Role{
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		CreatedBy:   createdBy,
		UpdatedAt:   time.Now(),
		UpdatedBy:   createdBy,
	}
}

// Update modifies the role's description and updates audit fields.
//
// Parameters:
//   - description: New description for the role
//   - updatedBy: Identifier of who is updating the role
func (r *Role) Update(description string, updatedBy string) {
	r.Description = description
	r.UpdatedAt = time.Now()
	r.UpdatedBy = updatedBy
}


