package user

import (
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
// Lists all users in the system
func ListUsers() []*User {
	queryResult, err := db.RunQuery(SELECT_ALL_USERS)
	if err != nil {
		return []*User{}
	}
	defer queryResult.Close()

	var users []*User
	for queryResult.Next() {
		var idStr, subject, passwordHash, salt, pepper, createdBy, updatedBy string
		var createdAt, updatedAt time.Time

		err := queryResult.Scan(&idStr, &subject, &passwordHash, &salt, &pepper, &createdAt, &createdBy, &updatedAt, &updatedBy)
		if err != nil {
			continue
		}

		id, err := uuid.Parse(idStr)
		if err != nil {
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

	return users
}

// Fetches a user by their ID
func GetUserById(userId uuid.UUID) *User {
	queryResult, err := db.RunQueryWithArgs(SELECT_USER_BY_ID, userId.String())
	if err != nil {
		return nil
	}
	defer queryResult.Close()

	if !queryResult.Next() {
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

	return &user
}

// Fetches a user by their subject
func GetUserBySubject(subject string) *User {
	queryResult, err := db.RunQueryWithArgs(SELECT_USER_BY_SUBJECT, subject)
	if err != nil {
		return nil
	}
	defer queryResult.Close()

	if !queryResult.Next() {
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

// Checks if a user exists by their subject
func ExistsUserBySubject(subject string) bool {
	user := GetUserBySubject(subject)
	return user != nil
}

// Persists the user to the database (insert or update), syncing roles and permissions as needed
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

// Synchronizes user roles between the provided user and the existing user on the database
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

// Synchronizes user permissions between the provided user and the existing user on the database
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

// Saves the user to the database
func SaveUser(user *User) {
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
		panic(err)
	}
}

// Updates the user in the database
func UpdateUser(user *User) {
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
		panic(err)
	}
}

// Deletes the user from the database
func DeleteUser(userId uuid.UUID) {
	err := db.RunCommandWithArgs(DELETE_USER, userId.String())
	if err != nil {
		panic(err)
	}
}

// Loads user roles from the database
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

// Adds a role to the user in the database
func AddUserRole(user *User, role UserRole) {
	err := db.RunCommandWithArgs(INSERT_USER_ROLE, user.Id.String(), string(role))
	if err != nil {
		panic(err)
	}
}

// Removes a role from the user in the database
func RemoveUserRole(user *User, role UserRole) {
	err := db.RunCommandWithArgs(DELETE_USER_ROLE, user.Id.String(), string(role))
	if err != nil {
		panic(err)
	}
}

// Loads user permissions from the database
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

// Adds a permission to the user in the database
func AddUserPermission(user *User, permission Permission) {
	err := db.RunCommandWithArgs(INSERT_USER_PERMISSION, user.Id.String(), string(permission))
	if err != nil {
		panic(err)
	}
}

// Removes a permission from the user in the database
func RemoveUserPermission(user *User, permission Permission) {
	err := db.RunCommandWithArgs(DELETE_USER_PERMISSION, user.Id.String(), string(permission))
	if err != nil {
		panic(err)
	}
}