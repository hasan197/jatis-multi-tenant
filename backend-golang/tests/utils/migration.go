package utils

import (
	"database/sql"
	"fmt"
)

// MigrateTestDB menjalankan migrasi untuk database test
func MigrateTestDB(db *sql.DB) error {
	// Buat tabel users
	query := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
	`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating users table: %v", err)
	}

	return nil
}

// CleanupTestDB membersihkan data test
func CleanupTestDB(db *sql.DB) error {
	_, err := db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")
	if err != nil {
		return fmt.Errorf("error cleaning up test data: %v", err)
	}
	return nil
} 