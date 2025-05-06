package http

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"sample-stack-golang/pkg/logger"
)

// GetDLQStatus handles getting dead-letter queue status for a tenant
func (h *TenantHandler) GetDLQStatus(c echo.Context) error {
	tenantID := c.Param("id")
	if tenantID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "tenant ID is required"})
	}

	logger.Log.WithFields(map[string]interface{}{
		"tenant_id": tenantID,
	}).Info("[DLQ] Checking DLQ status")

	// Get channel from RabbitMQ
	ch, err := h.tenantUseCase.GetChannel()
	if err != nil {
		logger.Log.WithFields(map[string]interface{}{
			"tenant_id": tenantID,
			"error":     err,
		}).Error("[DLQ] Failed to get channel")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get channel"})
	}
	defer ch.Close()

	// Check DLQ
	dlqName := fmt.Sprintf("dlq.tenant.%s", tenantID)
	dlq, err := ch.QueueInspect(dlqName)
	if err != nil {
		logger.Log.WithFields(map[string]interface{}{
			"tenant_id": tenantID,
			"queue":     dlqName,
			"error":     err,
		}).Warn("[DLQ] Failed to inspect DLQ, it may not exist yet")
		
		// DLQ mungkin belum ada, bukan error
		return c.JSON(http.StatusOK, map[string]interface{}{
			"exists":        false,
			"messageCount":  0,
			"consumerCount": 0,
		})
	}

	logger.Log.WithFields(map[string]interface{}{
		"tenant_id":     tenantID,
		"queue":         dlqName,
		"message_count": dlq.Messages,
		"consumer_count": dlq.Consumers,
	}).Info("[DLQ] Successfully checked DLQ status")

	return c.JSON(http.StatusOK, map[string]interface{}{
		"exists":        true,
		"messageCount":  dlq.Messages,
		"consumerCount": dlq.Consumers,
	})
}

// ActivateConsumer handles activating a consumer for a tenant
func (h *TenantHandler) ActivateConsumer(c echo.Context) error {
	tenantID := c.Param("id")
	if tenantID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "tenant ID is required"})
	}

	logger.Log.WithFields(map[string]interface{}{
		"tenant_id": tenantID,
	}).Info("Activating consumer for tenant")

	err := h.tenantUseCase.StartConsumer(c.Request().Context(), tenantID)
	if err != nil {
		logger.Log.WithFields(map[string]interface{}{
			"tenant_id": tenantID,
			"error":     err,
		}).Error("Failed to start consumer")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to start consumer: %v", err)})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Consumer activated successfully"})
}
