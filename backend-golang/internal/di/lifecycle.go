package di

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"sample-stack-golang/internal/config"
)

// LifecycleManager mengatur lifecycle dari aplikasi
type LifecycleManager struct {
	container *Container
	config    ConfigContainer
}

// NewLifecycleManager membuat instance baru dari LifecycleManager
func NewLifecycleManager(container *Container, config ConfigContainer) *LifecycleManager {
	return &LifecycleManager{
		container: container,
		config:    config,
	}
}

// Start menginisialisasi dan menjalankan aplikasi
func (lm *LifecycleManager) Start(ctx context.Context) error {
	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for shutdown signal
	go func() {
		sig := <-sigChan
		log.Printf("Received signal: %v", sig)
		if err := lm.Shutdown(ctx); err != nil {
			log.Printf("Error during shutdown: %v", err)
		}
		os.Exit(0)
	}()

	return nil
}

// Shutdown melakukan cleanup dan menutup semua resources
func (lm *LifecycleManager) Shutdown(ctx context.Context) error {
	// Close database connection
	if db := lm.config.GetDB(); db != nil {
		db.Close()
	}

	// Close container
	if err := lm.container.Close(); err != nil {
		return fmt.Errorf("error closing container: %w", err)
	}

	return nil
}

// Initialize menginisialisasi semua dependencies
func Initialize() (*LifecycleManager, error) {
	// Initialize configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize config: %w", err)
	}

	// Create DI container
	container := NewContainer()

	// Initialize database
	pool, err := initDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Create config container
	configContainer := NewConfigContainer(cfg, pool)

	// Register dependencies
	container.Register("config", cfg)
	container.Register("db", pool)
	container.RegisterCloser(func() error {
		if pool != nil {
			pool.Close()
		}
		return nil
	})

	// Create lifecycle manager
	manager := NewLifecycleManager(container, configContainer)

	return manager, nil
} 