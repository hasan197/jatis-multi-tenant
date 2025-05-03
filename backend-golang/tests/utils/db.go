package utils

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// GetTestDB membuat koneksi database untuk testing
func GetTestDB() (*sql.DB, error) {
	// Konfigurasi database test
	dbHost := "test-db"
	dbPort := 5432
	dbUser := "test"
	dbPassword := "test"
	dbName := "testdb"

	// Buat connection string
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	// Buka koneksi database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	// Test koneksi
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	return db, nil
} 