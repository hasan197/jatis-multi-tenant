package migration

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"sample-stack-golang/internal/config"
)

// Migration represents a database migration
type Migration struct {
	Version int
	Name    string
	Up      string
	Down    string
}

// RunMigrations menjalankan migrasi database
func RunMigrations(cfg *config.Config) error {
	// Build connection string
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Name,
	)

	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Create migrations table if not exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get applied migrations
	rows, err := db.Query("SELECT name FROM migrations ORDER BY id")
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("failed to scan migration name: %w", err)
		}
		applied[name] = true
	}

	// Run migrations
	for _, migration := range migrations {
		if applied[migration.Name] {
			continue
		}

		log.Printf("Running migration: %s", migration.Name)

		// Start transaction
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}

		// Run migration
		if _, err := tx.Exec(migration.Up); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to run migration %s: %w", migration.Name, err)
		}

		// Record migration
		if _, err := tx.Exec("INSERT INTO migrations (name) VALUES ($1)", migration.Name); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %w", migration.Name, err)
		}

		// Commit transaction
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", migration.Name, err)
		}

		log.Printf("Completed migration: %s", migration.Name)
	}

	return nil
}

// RollbackMigrations melakukan rollback migrasi database
func RollbackMigrations(cfg *config.Config) error {
	// Build database URL
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
	)

	// Open database connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := sql.Open("postgres", dbURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer pool.Close()

	// Rollback migrations
	if err := rollbackMigrations(ctx, pool); err != nil {
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}

	log.Println("Migrations rolled back successfully")
	return nil
}

// runMigrations menjalankan migrasi menggunakan pgx
func runMigrations(ctx context.Context, pool *sql.DB) error {
	// Get all migration files
	migrations, err := getMigrationFiles("migrations")
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	// Get current version
	currentVersion, err := getCurrentVersion(ctx, pool)
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	// Run pending migrations
	for _, migration := range migrations {
		if migration.Version > currentVersion {
			if err := runMigration(ctx, pool, migration); err != nil {
				return fmt.Errorf("failed to run migration %d: %w", migration.Version, err)
			}
		}
	}

	return nil
}

// rollbackMigrations melakukan rollback migrasi menggunakan pgx
func rollbackMigrations(ctx context.Context, pool *sql.DB) error {
	// Get current version
	currentVersion, err := getCurrentVersion(ctx, pool)
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	if currentVersion == 0 {
		return nil
	}

	// Get all migration files
	migrations, err := getMigrationFiles("migrations")
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	// Find current migration
	var currentMigration *Migration
	for _, m := range migrations {
		if m.Version == currentVersion {
			currentMigration = m
			break
		}
	}

	if currentMigration == nil {
		return fmt.Errorf("migration version %d not found", currentVersion)
	}

	// Rollback current migration
	if err := rollbackMigration(ctx, pool, currentMigration); err != nil {
		return fmt.Errorf("failed to rollback migration %d: %w", currentVersion, err)
	}

	return nil
}

// createMigrationsTable membuat tabel migrations jika belum ada
func createMigrationsTable(ctx context.Context, pool *sql.DB) error {
	_, err := pool.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	return err
}

// getCurrentVersion mendapatkan versi migrasi terakhir
func getCurrentVersion(ctx context.Context, pool *sql.DB) (int, error) {
	var version int
	err := pool.QueryRowContext(ctx, "SELECT COALESCE(MAX(version), 0) FROM migrations").Scan(&version)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return version, nil
}

// getMigrationFiles mendapatkan semua file migrasi
func getMigrationFiles(dir string) ([]*Migration, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var migrations []*Migration
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			parts := strings.Split(file.Name(), "_")
			if len(parts) < 2 {
				continue
			}

			version, err := strconv.Atoi(parts[0])
			if err != nil {
				continue
			}

			content, err := os.ReadFile(filepath.Join(dir, file.Name()))
			if err != nil {
				return nil, err
			}

			// Split content into up and down migrations
			parts := strings.Split(string(content), "-- +migrate Down")
			if len(parts) != 2 {
				continue
			}

			up := strings.TrimPrefix(parts[0], "-- +migrate Up\n")
			down := strings.TrimSpace(parts[1])

			migrations = append(migrations, &Migration{
				Version: version,
				Name:    strings.TrimSuffix(file.Name(), ".sql"),
				Up:      up,
				Down:    down,
			})
		}
	}

	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

// runMigration menjalankan satu migrasi
func runMigration(ctx context.Context, pool *sql.DB, migration *Migration) error {
	// Start transaction
	tx, err := pool.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Run migration
	if _, err := tx.ExecContext(ctx, migration.Up); err != nil {
		return err
	}

	// Record migration
	if _, err := tx.ExecContext(ctx, "INSERT INTO migrations (name) VALUES ($1)", migration.Name); err != nil {
		return err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// rollbackMigration melakukan rollback satu migrasi
func rollbackMigration(ctx context.Context, pool *sql.DB, migration *Migration) error {
	// Start transaction
	tx, err := pool.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Run rollback
	if _, err := tx.ExecContext(ctx, migration.Down); err != nil {
		return err
	}

	// Remove migration record
	if _, err := tx.ExecContext(ctx, "DELETE FROM migrations WHERE name = $1", migration.Name); err != nil {
		return err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
} 