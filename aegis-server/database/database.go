// Package database provides low-level database operations for SQLite.
// Handles connection management and query execution with automatic connection cleanup.
package database

import (
	"database/sql"
	"log"
	"os"
	_ "github.com/mattn/go-sqlite3"
)

var DB_FILE = getDBFile()

// getDBFile returns the database file path based on environment or testing flag.
// Returns "aegis-test.db" if AEGIS_TEST_MODE is set, otherwise uses AEGIS_DB_PATH env var or defaults to "/app/data/aegis.db".
func getDBFile() string {
	if os.Getenv("AEGIS_TEST_MODE") == "true" {
		return "aegis-test.db"
	}
	
	// Check for custom database path from environment
	dbPath := os.Getenv("AEGIS_DB_PATH")
	if dbPath != "" {
		return dbPath
	}
	
	// Default to /app/data/aegis.db for Docker persistence
	return "/app/data/aegis.db"
}

// SetTestMode enables test mode, using aegis-test.db instead of aegis.db
func SetTestMode() {
	DB_FILE = "aegis-test.db"
}

// OpenConnection establishes a new connection to the SQLite database.
// The connection should be closed by the caller using defer db.Close().
//
// Returns:
//   - *sql.DB: Database connection handle
//   - error: Error if connection fails
func OpenConnection() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", DB_FILE)
	if err != nil {
		log.Println("Failed to open database:", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Println("Failed to connect to database:", err)
		return nil, err
	}

	return db, nil
}

// RunCommand executes a SQL command (INSERT, UPDATE, DELETE, CREATE, etc.) without parameters.
// Opens and closes the database connection automatically.
//
// Parameters:
//   - query: The SQL command to execute
//
// Returns:
//   - error: Error if execution fails
func RunCommand(query string) error {
	db, err := OpenConnection();
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// RunCommandWithArgs executes a SQL command with parameterized arguments.
// Opens and closes the database connection automatically. Use this to prevent SQL injection.
//
// Parameters:
//   - query: The SQL command with ? placeholders
//   - args: Values to substitute for placeholders
//
// Returns:
//   - error: Error if execution fails
func RunCommandWithArgs(query string, args ...interface{}) error {
	db, err := OpenConnection();
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

// RunQuery executes a SQL SELECT query without parameters.
// Opens and closes the database connection automatically.
// The caller must close the returned rows using defer rows.Close().
//
// Parameters:
//   - query: The SQL SELECT query to execute
//
// Returns:
//   - *sql.Rows: Result set from the query
//   - error: Error if execution fails
func RunQuery(query string) (*sql.Rows, error) {
	db, err := OpenConnection();
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

// RunQueryWithArgs executes a SQL SELECT query with parameterized arguments.
// Opens and closes the database connection automatically. Use this to prevent SQL injection.
// The caller must close the returned rows using defer rows.Close().
//
// Parameters:
//   - query: The SQL SELECT query with ? placeholders
//   - args: Values to substitute for placeholders
//
// Returns:
//   - *sql.Rows: Result set from the query
//   - error: Error if execution fails
func RunQueryWithArgs(query string, args ...interface{}) (*sql.Rows, error) {
	db, err := OpenConnection();
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}