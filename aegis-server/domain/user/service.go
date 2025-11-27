package user

import (
	"log"
	"time"
	"github.com/google/uuid"
	db "nfcunha/aegis/database"
)

const ( 
	
	SELECT_ALL_USERS = `
		SELECT 
			id, 
			subject, 
			password_hash, 
			salt, 
			pepper, 
			created_at, 
			created_by, 
			updated_at, 
			updated_by 
		FROM 
			users
	`

	SELECT_USER_BY_ID = `
		SELECT 
			id, 
			subject, 
			password_hash, 
			salt, 
			pepper, 
			created_at, 
			created_by, 
			updated_at, 
			updated_by 
		FROM 
			users 
		WHERE 
			id = ?
	`

	SELECT_USER_BY_SUBJECT = `
		SELECT 
			id, 
			subject, 
			password_hash, 
			salt, 
			pepper, 
			created_at, 
			created_by, 
			updated_at, 
			updated_by 
		FROM 
			users 
		WHERE 
			subject = ?
	`

	SELECT_USER_ROLES = `
		SELECT 
			role 
		FROM 
			user_roles 
		WHERE 
			user_id = ?
	`

	SELECT_USER_PERMISSIONS = `
		SELECT 
			permission 
		FROM 
			user_permissions 
		WHERE 
			user_id = ?
	`

	INSERT_USER = `
		INSERT INTO users (
			id, 
			subject, 
			password_hash, 
			salt, 
			pepper, 
			created_at, 
			created_by, 
			updated_at, 
			updated_by
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	DELETE_USER = `
		DELETE FROM users 
		WHERE id = ?
	`

	UPDATE_USER = `
		UPDATE 
			users 
		SET 
			subject = ?, 
			password_hash = ?, 
			salt = ?, 
			pepper = ?, 
			updated_at = ?, 
			updated_by = ? 
		WHERE id = ?
	`

	INSERT_USER_ROLE = `
		INSERT INTO user_roles (user_id, role) 
		VALUES (?, ?)
	`

	DELETE_USER_ROLE = `
		DELETE FROM user_roles 
		WHERE user_id = ? AND role = ?
	`

	INSERT_USER_PERMISSION = `
		INSERT INTO user_permissions (user_id, permission) 
		VALUES (?, ?)
	`

	DELETE_USER_PERMISSION = `
		DELETE FROM user_permissions 
		WHERE user_id = ? AND permission = ?
	`
)

// ListUsers retrieves all users from the database including their roles and permissions.
//
// Returns:
//   - Slice of User pointers, empty slice if no users exist or on error
func ListUsers() []*User {
	log.Println("Listing all users")
	queryResult, err := db.RunQuery(SELECT_ALL_USERS)
	if err != nil {
		log.Println("Error listing users:", err)
		return []*User{}
	}
	defer queryResult.Close()

	var users []*User
	for queryResult.Next() {
		var idStr, subject, passwordHash, salt, pepper, createdBy, updatedBy string
		var createdAt, updatedAt time.Time

		err := queryResult.Scan(&idStr, &subject, &passwordHash, &salt, &pepper, &createdAt, &createdBy, &updatedAt, &updatedBy)
		if err != nil {
			log.Println("Error scanning user:", err)
			continue
		}

		id, err := uuid.Parse(idStr)
		if err != nil {
			log.Println("Error parsing user ID:", err)
			continue
		}

		user := &User{
			Id:           id,
			Subject:      subject,
			PasswordHash: passwordHash,
			Salt:         salt,
			Pepper:       pepper,
			CreatedAt:    createdAt,
			CreatedBy:    createdBy,
			UpdatedAt:    updatedAt,
			UpdatedBy:    updatedBy,
		}
		LoadUserPermissions(user)
		LoadUserRoles(user)
		users = append(users, user)
	}

	log.Printf("Found %d users", len(users))
	return users
}

// GetUserById retrieves a user by their unique identifier.
// Loads associated roles and permissions from the database.
//
// Parameters:
//   - userId: The UUID of the user to retrieve
//
// Returns:
//   - Pointer to the User if found, nil otherwise
func GetUserById(userId uuid.UUID) *User {
	log.Printf("Fetching user by ID: %s", userId.String())
	queryResult, err := db.RunQueryWithArgs(SELECT_USER_BY_ID, userId.String())
	if err != nil {
		log.Println("Error fetching user:", err)
		return nil
	}
	defer queryResult.Close()

	if !queryResult.Next() {
		log.Printf("User not found: %s", userId.String())
		return nil
	}

	var idStr, subject, passwordHash, salt, pepper, createdBy, updatedBy string
	var createdAt, updatedAt time.Time

	err = queryResult.Scan(&idStr, &subject, &passwordHash, &salt, &pepper, &createdAt, &createdBy, &updatedAt, &updatedBy)
	if err != nil {
		return nil
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil
	}

	user := User{
		Id:           id,
		Subject:      subject,
		PasswordHash: passwordHash,
		Salt:         salt,
		Pepper:       pepper,
		CreatedAt:    createdAt,
		CreatedBy:    createdBy,
		UpdatedAt:    updatedAt,
		UpdatedBy:    updatedBy,
	}

	LoadUserPermissions(&user)
	LoadUserRoles(&user)

	log.Printf("User found: %s", user.Subject)
	return &user
}

// GetUserBySubject retrieves a user by their subject identifier.
// Loads associated roles and permissions from the database.
//
// Parameters:
//   - subject: The subject identifier (typically email or username)
//
// Returns:
//   - Pointer to the User if found, nil otherwise
func GetUserBySubject(subject string) *User {
	log.Printf("Fetching user by subject: %s", subject)
	queryResult, err := db.RunQueryWithArgs(SELECT_USER_BY_SUBJECT, subject)
	if err != nil {
		log.Println("Error fetching user:", err)
		return nil
	}
	defer queryResult.Close()

	if !queryResult.Next() {
		log.Printf("User not found: %s", subject)
		return nil
	}

	var idStr, passwordHash, salt, pepper, createdBy, updatedBy string
	var createdAt, updatedAt time.Time

	err = queryResult.Scan(&idStr, &subject, &passwordHash, &salt, &pepper, &createdAt, &createdBy, &updatedAt, &updatedBy)
	if err != nil {
		return nil
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil
	}

	user := User{
		Id:           id,
		Subject:      subject,
		PasswordHash: passwordHash,
		Salt:         salt,
		Pepper:       pepper,
		CreatedAt:    createdAt,
		CreatedBy:	createdBy,
		UpdatedAt:    updatedAt,
		UpdatedBy:    updatedBy,
	}

	LoadUserPermissions(&user)
	LoadUserRoles(&user)

	return &user
}

// ExistsUserBySubject checks if a user with the given subject exists in the database.
//
// Parameters:
//   - subject: The subject identifier to check
//
// Returns:
//   - true if a user with this subject exists, false otherwise
func ExistsUserBySubject(subject string) bool {
	user := GetUserBySubject(subject)
	return user != nil
}

// PersistUser saves or updates a user in the database.
// If the user doesn't exist, inserts a new record. If it exists, updates the record
// and synchronizes roles and permissions by removing those no longer assigned and
// adding new ones.
//
// Parameters:
//   - user: The user to persist
func PersistUser(user *User) {
	if user == nil {
		return
	}

	existingUser := GetUserById(user.Id)
	if existingUser == nil {
		SaveUser(user)
	} else {
		UpdateUser(user)
		syncRoles(user, existingUser)
		syncPermissions(user, existingUser)
	}

	for _, role := range user.Roles {
		AddUserRole(user, role)
	}
	for _, permission := range user.Permissions {
		AddUserPermission(user, permission)
	}
}

// syncRoles synchronizes user roles by removing any roles from the existing user
// that are not present in the updated user. Uses a map for O(1) lookup performance.
//
// Parameters:
//   - user: The user with updated roles
//   - existingUser: The current user state from the database
func syncRoles(user, existingUser *User) {
	newRoles := make(map[UserRole]bool)
	for _, role := range user.Roles {
		newRoles[role] = true
	}
	for _, role := range existingUser.Roles {
		if !newRoles[role] {
			RemoveUserRole(user, role)
		}
	}
}

// syncPermissions synchronizes user permissions by removing any permissions from the existing user
// that are not present in the updated user. Uses a map for O(1) lookup performance.
//
// Parameters:
//   - user: The user with updated permissions
//   - existingUser: The current user state from the database
func syncPermissions(user, existingUser *User) {
	newPermissions := make(map[Permission]bool)
	for _, permission := range user.Permissions {
		newPermissions[permission] = true
	}
	for _, permission := range existingUser.Permissions {
		if !newPermissions[permission] {
			RemoveUserPermission(user, permission)
		}
	}
}

// SaveUser inserts a new user record into the database.
//
// Parameters:
//   - user: The user to save
//
// Panics:
//   - If the database insertion fails
func SaveUser(user *User) {
	log.Printf("Saving user: %s", user.Subject)
	err := db.RunCommandWithArgs(INSERT_USER,
		user.Id.String(),
		user.Subject,
		user.PasswordHash,
		user.Salt,
		user.Pepper,
		user.CreatedAt,
		user.CreatedBy,
		user.UpdatedAt,
		user.UpdatedBy,
	)

	if err != nil {
		log.Printf("Error saving user %s: %v", user.Subject, err)
		panic(err)
	}
	log.Printf("User saved successfully: %s", user.Subject)
}

// UpdateUser updates an existing user record in the database.
//
// Parameters:
//   - user: The user with updated data
//
// Panics:
//   - If the database update fails
func UpdateUser(user *User) {
	log.Printf("Updating user: %s", user.Subject)
	err := db.RunCommandWithArgs(UPDATE_USER,
		user.Subject,
		user.PasswordHash,
		user.Salt,
		user.Pepper,
		user.UpdatedAt,
		user.UpdatedBy,
		user.Id.String(),
	)

	if err != nil {
		log.Printf("Error updating user %s: %v", user.Subject, err)
		panic(err)
	}
	log.Printf("User updated successfully: %s", user.Subject)
}

// DeleteUser removes a user and all associated roles/permissions from the database.
// Foreign key constraints handle cascading deletes of roles and permissions.
//
// Parameters:
//   - userId: The UUID of the user to delete
//
// Panics:
//   - If the database deletion fails
func DeleteUser(userId uuid.UUID) {
	log.Printf("Deleting user: %s", userId.String())
	err := db.RunCommandWithArgs(DELETE_USER, userId.String())
	if err != nil {
		log.Printf("Error deleting user %s: %v", userId.String(), err)
		panic(err)
	}
	log.Printf("User deleted successfully: %s", userId.String())
}

// LoadUserRoles loads all roles assigned to a user from the database.
//
// Parameters:
//   - user: The user whose roles should be loaded
func LoadUserRoles(user *User) {
	rows, err := db.RunQueryWithArgs(SELECT_USER_ROLES, user.Id.String())
	if err != nil {
		return
	}
	defer rows.Close()

	var roles []UserRole
	for rows.Next() {
		var roleStr string
		err := rows.Scan(&roleStr)
		if err != nil {
			continue
		}
		role := UserRole(roleStr)
		roles = append(roles, role)
	}
	user.Roles = roles
}

// AddUserRole associates a role with a user in the database.
//
// Parameters:
//   - user: The user to add the role to
//   - role: The role to add
//
// Panics:
//   - If the database insertion fails
func AddUserRole(user *User, role UserRole) {
	err := db.RunCommandWithArgs(INSERT_USER_ROLE, user.Id.String(), string(role))
	if err != nil {
		panic(err)
	}
}

// RemoveUserRole removes a role association from a user in the database.
//
// Parameters:
//   - user: The user to remove the role from
//   - role: The role to remove
//
// Panics:
//   - If the database deletion fails
func RemoveUserRole(user *User, role UserRole) {
	err := db.RunCommandWithArgs(DELETE_USER_ROLE, user.Id.String(), string(role))
	if err != nil {
		panic(err)
	}
}

// LoadUserPermissions loads all permissions assigned to a user from the database.
//
// Parameters:
//   - user: The user whose permissions should be loaded
func LoadUserPermissions(user *User) {
	rows, err := db.RunQueryWithArgs(SELECT_USER_PERMISSIONS, user.Id.String())
	if err != nil {
		return
	}
	defer rows.Close()

	var permissions []Permission
	for rows.Next() {
		var permissionStr string
		err := rows.Scan(&permissionStr)
		if err != nil {
			continue
		}
		permission := Permission(permissionStr)
		permissions = append(permissions, permission)
	}
	user.Permissions = permissions
}

// AddUserPermission associates a permission with a user in the database.
//
// Parameters:
//   - user: The user to add the permission to
//   - permission: The permission to add
//
// Panics:
//   - If the database insertion fails
func AddUserPermission(user *User, permission Permission) {
	err := db.RunCommandWithArgs(INSERT_USER_PERMISSION, user.Id.String(), string(permission))
	if err != nil {
		panic(err)
	}
}

// RemoveUserPermission removes a permission association from a user in the database.
//
// Parameters:
//   - user: The user to remove the permission from
//   - permission: The permission to remove
//
// Panics:
//   - If the database deletion fails
func RemoveUserPermission(user *User, permission Permission) {
	err := db.RunCommandWithArgs(DELETE_USER_PERMISSION, user.Id.String(), string(permission))
	if err != nil {
		panic(err)
	}
}