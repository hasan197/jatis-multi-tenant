package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	
	"sample-stack-golang/internal/di"
	"sample-stack-golang/internal/metrics"
	userHttp "sample-stack-golang/internal/modules/user/delivery/http"
	"sample-stack-golang/pkg/logger"
)

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

	// Inisialisasi router Gin
	r := gin.Default()

	// Setup metrics
	metrics.SetupMetrics(r)

	// Konfigurasi CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"message":   "Server is running with hot reload!",
		})
	})

	// Endpoint hello-world
	r.GET("/api/hello-world", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
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
	userHttp.RegisterRoutes(r, userHandler)
	fmt.Println("User routes registered successfully")

	// Jalankan server
	fmt.Println("Starting server on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
} 