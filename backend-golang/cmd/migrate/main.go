package main

import (
	"flag"
	"log"
	"os"

	"sample-stack-golang/internal/config"
	"sample-stack-golang/internal/database/migration"
)

func main() {
	// Parse command line flags
	rollback := flag.Bool("rollback", false, "Rollback the last migration")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Run migrations
	if *rollback {
		if err := migration.RollbackMigrations(cfg.DatabaseURL); err != nil {
			log.Fatalf("Failed to rollback migrations: %v", err)
		}
		os.Exit(0)
	}

	if err := migration.RunMigrations(cfg.DatabaseURL); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
} 