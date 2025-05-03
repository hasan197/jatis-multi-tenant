package http

import (
	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers tenant routes
func RegisterRoutes(e *echo.Echo, handler *TenantHandler) {
	tenants := e.Group("/tenants")
	tenants.POST("", handler.Create)
	tenants.GET("", handler.List)
	tenants.GET("/:id", handler.GetByID)
	tenants.PUT("/:id", handler.Update)
	tenants.DELETE("/:id", handler.Delete)
} 