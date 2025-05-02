package di

import (
	"database/sql"
	"fmt"
	"sample-stack-golang/internal/config"
)

// ConfigContainer adalah interface untuk mengakses konfigurasi
type ConfigContainer interface {
	GetConfig() *config.Config
	GetDB() *sql.DB
}

type configContainer struct {
	config *config.Config
	db     *sql.DB
}

// NewConfigContainer membuat instance baru dari ConfigContainer
func NewConfigContainer(cfg *config.Config, db *sql.DB) ConfigContainer {
	return &configContainer{
		config: cfg,
		db:     db,
	}
}

// GetConfig mengambil konfigurasi dari container
func (cc *configContainer) GetConfig() *config.Config {
	return cc.config
}

// GetDB mengambil database connection dari container
func (cc *configContainer) GetDB() *sql.DB {
	return cc.db
}

// InitializeConfig menginisialisasi konfigurasi dan dependencies yang diperlukan
func InitializeConfig() (ConfigContainer, error) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize database connection
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test database connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return NewConfigContainer(cfg, db), nil
} 