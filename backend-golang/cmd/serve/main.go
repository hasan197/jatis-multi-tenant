package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/go-playground/validator/v10"
	_ "github.com/lib/pq"
	
	"sample-stack-golang/internal/di"
	"sample-stack-golang/internal/metrics"
	userHttp "sample-stack-golang/internal/modules/user/delivery/http"
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
	
	// Inisialisasi logger
	logConfig := &logger.Config{
		LogLevel:      "debug",
		LogFilePath:   "logs/app.log",
		MaxSize:       10,
		MaxBackups:    5,
		MaxAge:        30,
		Compress:      true,
		ConsoleOutput: true,
	}
	
	if err := logger.InitLogger(logConfig); err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	fmt.Println("Logger initialized successfully")
	
	fmt.Println("Initializing DI container...")

	// Inisialisasi DI container dan lifecycle manager
	manager, err := di.Initialize()
	if err != nil {
		log.Fatal("Failed to initialize DI container:", err)
	}
	fmt.Println("DI container initialized successfully")

	// Setup context dengan timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Start lifecycle manager
	fmt.Println("Starting lifecycle manager...")
	if err := manager.Start(ctx); err != nil {
		log.Fatal("Failed to start lifecycle manager:", err)
	}
	fmt.Println("Lifecycle manager started successfully")

	// Inisialisasi Echo
	e := echo.New()

	// Setup validator
	e.Validator = &CustomValidator{validator: validator.New()}

	// Setup metrics
	metrics.SetupMetrics(e)

	// Konfigurasi CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		ExposeHeaders:    []string{echo.HeaderContentLength},
		AllowCredentials: true,
	}))

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":    "ok",
			"message":   "Server is running with hot reload!",
		})
	})

	// Endpoint hello-world
	e.GET("/api/hello-world", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message":   "Hello World dari Backend Golang!",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})
	
	// Setup user handler menggunakan DI container
	fmt.Println("Setting up user handler...")
	userService := manager.Services().GetUserService()
	if userService == nil {
		log.Fatal("Failed to get user service from DI container")
	}
	fmt.Println("User service retrieved successfully")
	
	userHandler := userHttp.NewUserHandler(userService)
	fmt.Println("User handler created successfully")
	
	// Register user routes
	fmt.Println("Registering user routes...")
	userHttp.RegisterRoutes(e, userHandler)
	fmt.Println("User routes registered successfully")

	// Jalankan server
	fmt.Println("Starting server on :8080...")
	if err := e.Start(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
} 