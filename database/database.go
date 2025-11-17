package database

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

const DB_FILE = "aegis.db"

func OpenConnection() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", DB_FILE)
	if err != nil {
		fmt.Println("Failed to open database:", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		return nil, err
	}

	return db, nil
}

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