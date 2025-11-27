package role

import (
	"log"
	"time"
	db "nfcunha/aegis/database"
)

const (
	SELECT_ALL_ROLES = `
		SELECT 
			name, 
			description, 
			created_at, 
			created_by, 
			updated_at, 
			updated_by 
		FROM 
			roles
	`

	SELECT_ROLE_BY_NAME = `
		SELECT 
			name, 
			description, 
			created_at, 
			created_by, 
			updated_at, 
			updated_by 
		FROM 
			roles 
		WHERE 
			name = ?
	`

	INSERT_ROLE = `
		INSERT INTO roles (
			name, 
			description, 
			created_at, 
			created_by, 
			updated_at, 
			updated_by
		) VALUES (?, ?, ?, ?, ?, ?)
	`

	UPDATE_ROLE = `
		UPDATE 
			roles 
		SET 
			description = ?, 
			updated_at = ?, 
			updated_by = ? 
		WHERE name = ?
	`

	DELETE_ROLE = `
		DELETE FROM roles 
		WHERE name = ?
	`
)

// ListRoles retrieves all roles from the database.
//
// Returns:
//   - Slice of Role pointers, empty slice if no roles exist or on error
func ListRoles() []*Role {
	log.Println("Listing all roles")
	queryResult, err := db.RunQuery(SELECT_ALL_ROLES)
	if err != nil {
		log.Println("Error listing roles:", err)
		return []*Role{}
	}
	defer queryResult.Close()

	var roles []*Role
	for queryResult.Next() {
		var name, description, createdBy, updatedBy string
		var createdAt, updatedAt time.Time

		err := queryResult.Scan(&name, &description, &createdAt, &createdBy, &updatedAt, &updatedBy)
		if err != nil {
			log.Println("Error scanning role:", err)
			continue
		}

		role := &Role{
			Name:        name,
			Description: description,
			CreatedAt:   createdAt,
			CreatedBy:   createdBy,
			UpdatedAt:   updatedAt,
			UpdatedBy:   updatedBy,
		}
		roles = append(roles, role)
	}

	log.Printf("Found %d roles", len(roles))
	return roles
}

// GetRoleByName retrieves a role by its unique name identifier.
//
// Parameters:
//   - name: The name of the role to retrieve
//
// Returns:
//   - Pointer to the Role if found, nil otherwise
func GetRoleByName(name string) *Role {
	log.Printf("Fetching role by name: %s", name)
	queryResult, err := db.RunQueryWithArgs(SELECT_ROLE_BY_NAME, name)
	if err != nil {
		log.Println("Error fetching role:", err)
		return nil
	}
	defer queryResult.Close()

	if !queryResult.Next() {
		log.Printf("Role not found: %s", name)
		return nil
	}

	var description, createdBy, updatedBy string
	var createdAt, updatedAt time.Time

	err = queryResult.Scan(&name, &description, &createdAt, &createdBy, &updatedAt, &updatedBy)
	if err != nil {
		log.Println("Error scanning role:", err)
		return nil
	}

	log.Printf("Role found: %s", name)
	return &Role{
		Name:        name,
		Description: description,
		CreatedAt:   createdAt,
		CreatedBy:   createdBy,
		UpdatedAt:   updatedAt,
		UpdatedBy:   updatedBy,
	}
}

// ExistsRoleByName checks if a role with the given name exists in the database.
//
// Parameters:
//   - name: The name to check
//
// Returns:
//   - true if a role with this name exists, false otherwise
func ExistsRoleByName(name string) bool {
	role := GetRoleByName(name)
	return role != nil
}

// PersistRole saves or updates a role in the database.
// If the role doesn't exist, inserts a new record. Otherwise, updates the existing one.
//
// Parameters:
//   - role: The role to persist
func PersistRole(role *Role) {
	if role == nil {
		return
	}

	existingRole := GetRoleByName(role.Name)
	if existingRole == nil {
		SaveRole(role)
	} else {
		UpdateRoleData(role)
	}
}

// SaveRole inserts a new role record into the database.
//
// Parameters:
//   - role: The role to save
//
// Panics:
//   - If the database insertion fails
func SaveRole(role *Role) {
	log.Printf("Saving role: %s", role.Name)
	err := db.RunCommandWithArgs(INSERT_ROLE,
		role.Name,
		role.Description,
		role.CreatedAt,
		role.CreatedBy,
		role.UpdatedAt,
		role.UpdatedBy,
	)

	if err != nil {
		log.Printf("Error saving role %s: %v", role.Name, err)
		panic(err)
	}
	log.Printf("Role saved successfully: %s", role.Name)
}

// UpdateRoleData updates an existing role record in the database.
//
// Parameters:
//   - role: The role with updated data
//
// Panics:
//   - If the database update fails
func UpdateRoleData(role *Role) {
	log.Printf("Updating role: %s", role.Name)
	err := db.RunCommandWithArgs(UPDATE_ROLE,
		role.Description,
		role.UpdatedAt,
		role.UpdatedBy,
		role.Name,
	)

	if err != nil {
		log.Printf("Error updating role %s: %v", role.Name, err)
		panic(err)
	}
	log.Printf("Role updated successfully: %s", role.Name)
}

// DeleteRole removes a role from the database.
// Associated user-role relationships are automatically deleted via foreign key constraints.
//
// Parameters:
//   - name: The name of the role to delete
//
// Panics:
//   - If the database deletion fails
func DeleteRole(name string) {
	log.Printf("Deleting role: %s", name)
	err := db.RunCommandWithArgs(DELETE_ROLE, name)
	if err != nil {
		log.Printf("Error deleting role %s: %v", name, err)
		panic(err)
	}
	log.Printf("Role deleted successfully: %s", name)
}
