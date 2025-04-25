package http

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes mendaftarkan semua endpoint untuk user
func RegisterRoutes(router *gin.Engine, handler *UserHandler) {
	userRoutes := router.Group("/api/users")
	{
		userRoutes.GET("", handler.GetUsers)
		userRoutes.GET("/:id", handler.GetUser)
		userRoutes.POST("", handler.CreateUser)
		userRoutes.PUT("/:id", handler.UpdateUser)
		userRoutes.DELETE("/:id", handler.DeleteUser)
	}
} 