package http

import (
	"github.com/labstack/echo/v4"
)

// RegisterRoutes mendaftarkan semua route untuk user
func RegisterRoutes(e *echo.Echo, h *UserHandler) {
	// Group routes untuk user
	users := e.Group("/api/users")
	
	// Register semua endpoint
	users.GET("", h.GetUsers)
	users.GET("/:id", h.GetUser)
	users.POST("", h.CreateUser)
	users.PUT("/:id", h.UpdateUser)
	users.DELETE("/:id", h.DeleteUser)
} 