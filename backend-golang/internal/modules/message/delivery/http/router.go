package http

import (
	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers all message routes
func (h *MessageHandler) RegisterRoutes(e *echo.Echo) {
	messageGroup := e.Group("/api/v1/tenants/:tenant_id/messages")
	messageGroup.POST("", h.Create)
	messageGroup.GET("", h.GetByTenant)
	messageGroup.GET("/:id", h.GetByID)
	messageGroup.PUT("/:id", h.Update)
	messageGroup.DELETE("/:id", h.Delete)
} 