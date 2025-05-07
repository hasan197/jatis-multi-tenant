package di

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jatis/sample-stack-golang/internal/config"
)

// ConfigContainer adalah interface untuk mengakses konfigurasi
type ConfigContainer interface {
	GetConfig() *config.Config
	GetDB() *pgxpool.Pool
}

type configContainer struct {
	config *config.Config
	pool   *pgxpool.Pool
}

// NewConfigContainer membuat instance baru dari ConfigContainer
func NewConfigContainer(cfg *config.Config, pool *pgxpool.Pool) ConfigContainer {
	return &configContainer{
		config: cfg,
		pool:   pool,
	}
}

// GetConfig mengambil konfigurasi dari container
func (cc *configContainer) GetConfig() *config.Config {
	return cc.config
}

// GetDB mengambil database connection dari container
func (cc *configContainer) GetDB() *pgxpool.Pool {
	return cc.pool
}

// InitializeConfig menginisialisasi konfigurasi dan dependencies yang diperlukan
func InitializeConfig() (ConfigContainer, error) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize database connection
	pool, err := pgxpool.New(context.Background(), cfg.DB.DatabaseURL())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test database connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return NewConfigContainer(cfg, pool), nil
}
