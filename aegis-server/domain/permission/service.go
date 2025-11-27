package permission

import (
	"log"
	"time"
	db "nfcunha/aegis/database"
)

const (
	SELECT_ALL_PERMISSIONS = `
		SELECT 
			name, 
			description, 
			created_at, 
			created_by, 
			updated_at, 
			updated_by 
		FROM 
			permissions
	`

	SELECT_PERMISSION_BY_NAME = `
		SELECT 
			name, 
			description, 
			created_at, 
			created_by, 
			updated_at, 
			updated_by 
		FROM 
			permissions 
		WHERE 
			name = ?
	`

	INSERT_PERMISSION = `
		INSERT INTO permissions (
			name, 
			description, 
			created_at, 
			created_by, 
			updated_at, 
			updated_by
		) VALUES (?, ?, ?, ?, ?, ?)
	`

	UPDATE_PERMISSION = `
		UPDATE 
			permissions 
		SET 
			description = ?, 
			updated_at = ?, 
			updated_by = ? 
		WHERE name = ?
	`

	DELETE_PERMISSION = `
		DELETE FROM permissions 
		WHERE name = ?
	`
)

// ListPermissions retrieves all permissions from the database.
//
// Returns:
//   - Slice of Permission pointers, empty slice if no permissions exist or on error
func ListPermissions() []*Permission {
	log.Println("Listing all permissions")
	queryResult, err := db.RunQuery(SELECT_ALL_PERMISSIONS)
	if err != nil {
		log.Println("Error listing permissions:", err)
		return []*Permission{}
	}
	defer queryResult.Close()

	var permissions []*Permission
	for queryResult.Next() {
		var name, description, createdBy, updatedBy string
		var createdAt, updatedAt time.Time

		err := queryResult.Scan(&name, &description, &createdAt, &createdBy, &updatedAt, &updatedBy)
		if err != nil {
			log.Println("Error scanning permission:", err)
			continue
		}

		permission := &Permission{
			Name:        name,
			Description: description,
			CreatedAt:   createdAt,
			CreatedBy:   createdBy,
			UpdatedAt:   updatedAt,
			UpdatedBy:   updatedBy,
		}
		permissions = append(permissions, permission)
	}

	log.Printf("Found %d permissions", len(permissions))
	return permissions
}

// GetPermissionByName retrieves a permission by its unique name identifier.
//
// Parameters:
//   - name: The name of the permission to retrieve
//
// Returns:
//   - Pointer to the Permission if found, nil otherwise
func GetPermissionByName(name string) *Permission {
	log.Printf("Fetching permission by name: %s", name)
	queryResult, err := db.RunQueryWithArgs(SELECT_PERMISSION_BY_NAME, name)
	if err != nil {
		log.Println("Error fetching permission:", err)
		return nil
	}
	defer queryResult.Close()

	if !queryResult.Next() {
		log.Printf("Permission not found: %s", name)
		return nil
	}

	var description, createdBy, updatedBy string
	var createdAt, updatedAt time.Time

	err = queryResult.Scan(&name, &description, &createdAt, &createdBy, &updatedAt, &updatedBy)
	if err != nil {
		log.Println("Error scanning permission:", err)
		return nil
	}

	log.Printf("Permission found: %s", name)
	return &Permission{
		Name:        name,
		Description: description,
		CreatedAt:   createdAt,
		CreatedBy:   createdBy,
		UpdatedAt:   updatedAt,
		UpdatedBy:   updatedBy,
	}
}

// ExistsPermissionByName checks if a permission with the given name exists in the database.
//
// Parameters:
//   - name: The name to check
//
// Returns:
//   - true if a permission with this name exists, false otherwise
func ExistsPermissionByName(name string) bool {
	permission := GetPermissionByName(name)
	return permission != nil
}

// PersistPermission saves or updates a permission in the database.
// If the permission doesn't exist, inserts a new record. Otherwise, updates the existing one.
//
// Parameters:
//   - permission: The permission to persist
func PersistPermission(permission *Permission) {
	if permission == nil {
		return
	}

	existingPermission := GetPermissionByName(permission.Name)
	if existingPermission == nil {
		SavePermission(permission)
	} else {
		UpdatePermissionData(permission)
	}
}

// SavePermission inserts a new permission record into the database.
//
// Parameters:
//   - permission: The permission to save
//
// Panics:
//   - If the database insertion fails
func SavePermission(permission *Permission) {
	log.Printf("Saving permission: %s", permission.Name)
	err := db.RunCommandWithArgs(INSERT_PERMISSION,
		permission.Name,
		permission.Description,
		permission.CreatedAt,
		permission.CreatedBy,
		permission.UpdatedAt,
		permission.UpdatedBy,
	)

	if err != nil {
		log.Printf("Error saving permission %s: %v", permission.Name, err)
		panic(err)
	}
	log.Printf("Permission saved successfully: %s", permission.Name)
}

// UpdatePermissionData updates an existing permission record in the database.
//
// Parameters:
//   - permission: The permission with updated data
//
// Panics:
//   - If the database update fails
func UpdatePermissionData(permission *Permission) {
	log.Printf("Updating permission: %s", permission.Name)
	err := db.RunCommandWithArgs(UPDATE_PERMISSION,
		permission.Description,
		permission.UpdatedAt,
		permission.UpdatedBy,
		permission.Name,
	)

	if err != nil {
		log.Printf("Error updating permission %s: %v", permission.Name, err)
		panic(err)
	}
	log.Printf("Permission updated successfully: %s", permission.Name)
}

// DeletePermission removes a permission from the database.
// Associated user-permission relationships are automatically deleted via foreign key constraints.
//
// Parameters:
//   - name: The name of the permission to delete
//
// Panics:
//   - If the database deletion fails
func DeletePermission(name string) {
	log.Printf("Deleting permission: %s", name)
	err := db.RunCommandWithArgs(DELETE_PERMISSION, name)
	if err != nil {
		log.Printf("Error deleting permission %s: %v", name, err)
		panic(err)
	}
	log.Printf("Permission deleted successfully: %s", name)
}
