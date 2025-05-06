package http

import (
	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers tenant routes
func RegisterRoutes(e *echo.Echo, h *TenantHandler) {
	// Tenant routes
	tenants := e.Group("/api/tenants")
	tenants.POST("", h.CreateTenant)
	tenants.DELETE("/:id", h.DeleteTenant)
	tenants.GET("/consumers", h.GetTenantConsumers)
	tenants.GET("/:id/consumers", h.GetTenantConsumers)
	tenants.PUT("/:id/config/concurrency", h.UpdateConcurrency) // New endpoint for configuring concurrency
	
	// RabbitMQ Publisher endpoints
	tenants.POST("/:id/publish", h.PublishMessage)      // Endpoint for publishing messages to RabbitMQ
	tenants.GET("/:id/queue-status", h.GetQueueStatus) // Endpoint for getting queue status
	
	// tenants.POST("", h.Create)
	tenants.GET("", h.List)
	tenants.GET("/:id", h.GetByID)
	tenants.PUT("/:id", h.Update)
	tenants.DELETE("/:id", h.Delete)
}