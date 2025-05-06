package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/go-playground/validator/v10"
	_ "github.com/lib/pq"
	
	"sample-stack-golang/internal/config"
	"sample-stack-golang/internal/di"
	"sample-stack-golang/pkg/infrastructure/metrics"
	"sample-stack-golang/pkg/graceful"
	userHttp "sample-stack-golang/internal/modules/user/delivery/http"
	tenantHttp "sample-stack-golang/internal/modules/tenant/delivery/http"
	messageHttp "sample-stack-golang/internal/modules/message/delivery/http"
	"sample-stack-golang/pkg/logger"
)

// CustomValidator adalah custom validator untuk Echo
type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	fmt.Println("Starting application with hot reload...")
	
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	fmt.Println("Configuration loaded successfully")
	
	// Inisialisasi logger
	if err := logger.InitLogger(cfg); err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	fmt.Println("Logger initialized successfully")
	
	fmt.Println("Initializing DI container...")

	// Initialize service
	service, err := di.NewService(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize service: %v", err)
	}
	defer service.Close()

	// Initialize Echo
	e := echo.New()

	// Setup validator
	e.Validator = &CustomValidator{validator: validator.New()}

	// Setup metrics
	metrics.SetupMetrics(e)

	// Create shutdown manager
	shutdownManager := graceful.NewShutdownManager(e, service.Close)

	// Set shutdown manager for tenant manager if available
	if tenantManager, ok := service.TenantUseCase.(interface{ SetShutdownManager(*graceful.ShutdownManager) }); ok {
		tenantManager.SetShutdownManager(shutdownManager)
	}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(shutdownManager.WaitGroupMiddleware())

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":    "ok",
			"message":   "Server is running with hot reload!",
			"version":   cfg.App.Version,
			"env":       cfg.App.Env,
		})
	})

	// Endpoint hello-world
	e.GET("/api/hello-world", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message":   "Hello World dari Backend Golang!",
			"timestamp": time.Now().Format(time.RFC3339),
			"version":   cfg.App.Version,
		})
	})
	
	// Initialize handlers
	userHandler := userHttp.NewUserHandler(service.UserUseCase)
	tenantHandler := tenantHttp.NewTenantHandler(service.TenantUseCase)
	messageHandler := messageHttp.NewMessageHandler(service.MessageUseCase)

	// Register routes
	userHttp.RegisterRoutes(e, userHandler)
	tenantHttp.RegisterRoutes(e, tenantHandler)
	messageHandler.RegisterRoutes(e)

	// Start server in a goroutine
	go func() {
		fmt.Printf("Server configuration - Port: %d\n", cfg.Server.Port)
		fmt.Printf("Starting server on port: %d\n", cfg.Server.Port)
		if err := e.Start(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil && err != http.ErrServerClosed {
			log.Printf("Failed to start server: %v", err)
			if closeErr := service.Close(); closeErr != nil {
				log.Printf("Failed to close service after server error: %v", closeErr)
			}
			os.Exit(1)
		}
	}()

	// Wait for graceful shutdown (this handles everything: signal capture, server shutdown, waiting for active processes)
	shutdownManager.WaitForShutdown()

	fmt.Println("Server shutdown complete")
} 