package database

// Migrate creates the database schema if it doesn't already exist.
// Creates five tables: users, roles, permissions, user_roles, and user_permissions.
// Includes foreign key constraints with CASCADE delete for referential integrity.
// This function is idempotent and safe to call multiple times.
func Migrate() {
	RunCommand(`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			subject TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			salt TEXT NOT NULL,
			pepper TEXT NOT NULL,
			created_at DATETIME NOT NULL,
			created_by TEXT NOT NULL,
			updated_at DATETIME NOT NULL,
			updated_by TEXT NOT NULL
	)`)
	RunCommand(`
		CREATE TABLE IF NOT EXISTS roles (
			name TEXT PRIMARY KEY,
			description TEXT,
			created_at DATETIME NOT NULL,
			created_by TEXT NOT NULL,
			updated_at DATETIME NOT NULL,
			updated_by TEXT NOT NULL
	)`)
	RunCommand(`
		CREATE TABLE IF NOT EXISTS permissions (
			name TEXT PRIMARY KEY,
			description TEXT,
			created_at DATETIME NOT NULL,
			created_by TEXT NOT NULL,
			updated_at DATETIME NOT NULL,
			updated_by TEXT NOT NULL
	)`)
	RunCommand(`
		CREATE TABLE IF NOT EXISTS user_roles (
			user_id TEXT NOT NULL,
			role TEXT NOT NULL,
			PRIMARY KEY (user_id, role),
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (role) REFERENCES roles(name) ON DELETE CASCADE
	)`)
	RunCommand(`
		CREATE TABLE IF NOT EXISTS user_permissions (
			user_id TEXT NOT NULL,
			permission TEXT NOT NULL,
			PRIMARY KEY (user_id, permission),
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (permission) REFERENCES permissions(name) ON DELETE CASCADE
	)`)
}