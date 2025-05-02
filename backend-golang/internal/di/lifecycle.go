package di

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// LifecycleManager mengatur lifecycle dari aplikasi
type LifecycleManager struct {
	container *Container
	config    ConfigContainer
	services  ServiceContainer
}

// NewLifecycleManager membuat instance baru dari LifecycleManager
func NewLifecycleManager(container *Container, config ConfigContainer, services ServiceContainer) *LifecycleManager {
	return &LifecycleManager{
		container: container,
		config:    config,
		services:  services,
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
		if err := db.Close(); err != nil {
			return fmt.Errorf("error closing database: %w", err)
		}
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
	configContainer, err := InitializeConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize config: %w", err)
	}

	// Create DI container
	container := NewContainer()

	// Register services
	RegisterServices(container, configContainer.GetDB(), configContainer.GetConfig())

	// Create service container
	serviceContainer := NewServiceContainer(container)

	// Create lifecycle manager
	manager := NewLifecycleManager(container, configContainer, serviceContainer)

	return manager, nil
}

// Services mengambil ServiceContainer
func (lm *LifecycleManager) Services() ServiceContainer {
	return lm.services
} 