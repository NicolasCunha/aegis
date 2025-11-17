package database

func Migrate() {
	RunCommand(`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			subject TEXT NOT NULL,
			password_hash TEXT NOT NULL,
			salt TEXT NOT NULL,
			pepper TEXT NOT NULL,
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
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	)`)
	RunCommand(`
		CREATE TABLE IF NOT EXISTS user_permissions (
			user_id TEXT NOT NULL,
			permission TEXT NOT NULL,
			PRIMARY KEY (user_id, permission),
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	)`)
}