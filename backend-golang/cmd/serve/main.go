package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	
	userHttp "sample-stack/internal/modules/user/delivery/http"
	userRepo "sample-stack/internal/modules/user/repository/postgresql"
	userUseCase "sample-stack/internal/modules/user/usecase"
)

func main() {
	// Inisialisasi koneksi database
	db, err := sql.Open("postgres", "postgres://postgres:postgres@postgres:5432/sample_db?sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()
	
	// Test koneksi database
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	
	// Inisialisasi router Gin
	r := gin.Default()

	// Konfigurasi CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Endpoint hello-world
	r.GET("/api/hello-world", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":   "Hello World dari Backend Golang!",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})
	
	// Setup dependencies untuk user module
	userRepository := userRepo.NewUserRepository(db)
	userUsecase := userUseCase.NewUserUseCase(userRepository)
	userHandler := userHttp.NewUserHandler(userUsecase)
	
	// Register user routes
	userHttp.RegisterRoutes(r, userHandler)

	// Jalankan server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
} 