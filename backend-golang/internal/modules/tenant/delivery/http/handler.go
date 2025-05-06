package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/streadway/amqp"
	"sample-stack-golang/internal/modules/tenant/domain"
	"sample-stack-golang/pkg/infrastructure/metrics"
	"sample-stack-golang/pkg/logger"
)

// TenantHandler handles HTTP requests for tenants
type TenantHandler struct {
	tenantUseCase domain.TenantUseCase
}

// NewTenantHandler creates a new tenant handler
func NewTenantHandler(tenantUseCase domain.TenantUseCase) *TenantHandler {
	return &TenantHandler{
		tenantUseCase: tenantUseCase,
	}
}

// Create handles tenant creation
func (h *TenantHandler) Create(c echo.Context) error {
	var tenant domain.Tenant
	if err := c.Bind(&tenant); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	tenant.CreatedAt = time.Now()
	tenant.UpdatedAt = time.Now()

	if err := h.tenantUseCase.Create(c.Request().Context(), &tenant); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, tenant)
}

// GetByID handles getting a tenant by ID
func (h *TenantHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	tenant, err := h.tenantUseCase.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, tenant)
}

// Update handles tenant updates
func (h *TenantHandler) Update(c echo.Context) error {
	id := c.Param("id")
	var tenant domain.Tenant
	if err := c.Bind(&tenant); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	tenant.ID = id
	tenant.UpdatedAt = time.Now()

	if err := h.tenantUseCase.Update(c.Request().Context(), &tenant); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, tenant)
}

// Delete handles tenant deletion
func (h *TenantHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.tenantUseCase.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// List handles listing all tenants
func (h *TenantHandler) List(c echo.Context) error {
	tenants, err := h.tenantUseCase.List(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, tenants)
}

// CreateTenant handles tenant creation with consumer
func (h *TenantHandler) CreateTenant(c echo.Context) error {
	var tenant domain.Tenant
	if err := c.Bind(&tenant); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	tenant.CreatedAt = time.Now()
	tenant.UpdatedAt = time.Now()

	if err := h.tenantUseCase.Create(c.Request().Context(), &tenant); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Start consumer for new tenant
	if err := h.tenantUseCase.StartConsumer(c.Request().Context(), tenant.ID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, tenant)
}

// DeleteTenant handles tenant deletion with consumer cleanup
func (h *TenantHandler) DeleteTenant(c echo.Context) error {
	id := c.Param("id")
	ctx := c.Request().Context()

	// Get tenant details for logging
	tenant, err := h.tenantUseCase.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "tenant not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Tenant not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Step 1: Stop consumer first
	if err := h.tenantUseCase.StopConsumer(ctx, id); err != nil {
		// Log the error but continue with deletion
		logger.Log.WithFields(map[string]interface{}{
			"tenant_id": id,
			"error":     err.Error(),
		}).Warn("Failed to stop consumer, proceeding with tenant deletion")
	}

	// Step 2: Delete tenant (this will also drop the message partition)
	if err := h.tenantUseCase.Delete(ctx, id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Log successful deletion
	logger.Log.WithFields(map[string]interface{}{
		"tenant_id":   id,
		"tenant_name": tenant.Name,
	}).Info("Tenant successfully deleted with all resources cleaned up")

	return c.NoContent(http.StatusNoContent)
}

// GetTenantConsumers handles getting all tenant consumers
func (h *TenantHandler) GetTenantConsumers(c echo.Context) error {
	// Get tenant ID from path parameter
	tenantID := c.Param("id")
	if tenantID != "" {
		// Get consumer for specific tenant
		consumers, err := h.tenantUseCase.GetConsumers(c.Request().Context())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		// Filter consumers by tenant ID
		for _, consumer := range consumers {
			if consumer.TenantID == tenantID {
				return c.JSON(http.StatusOK, consumer)
			}
		}

		return c.JSON(http.StatusNotFound, map[string]string{"error": "consumer not found"})
	}

	// Get all consumers
	consumers, err := h.tenantUseCase.GetConsumers(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, consumers)
}

// UpdateConcurrency handles updating tenant concurrency configuration
func (h *TenantHandler) UpdateConcurrency(c echo.Context) error {
	// Get tenant ID from path parameter
	id := c.Param("id")

	// Parse request body
	var config domain.ConcurrencyConfig
	if err := c.Bind(&config); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}

	// Validate worker count
	if config.Workers <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Worker count must be greater than 0"})
	}

	// Update concurrency configuration
	if err := h.tenantUseCase.UpdateConcurrency(c.Request().Context(), id, &config); err != nil {
		if err.Error() == "tenant not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Tenant not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Concurrency configuration updated successfully",
		"tenant_id": id,
		"workers": config.Workers,
	})
}

// GetQueueStatus handles getting queue status for a tenant
func (h *TenantHandler) GetQueueStatus(c echo.Context) error {
	tenantID := c.Param("id")
	if tenantID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "tenant ID is required"})
	}

	// Get consumer from manager
	consumer := h.tenantUseCase.GetConsumer(tenantID)

	// Get channel from RabbitMQ
	ch, err := h.tenantUseCase.GetChannel()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get channel"})
	}
	defer ch.Close()

	queueName := fmt.Sprintf("tenant.%s", tenantID)
	queue, err := ch.QueueInspect(queueName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to inspect queue"})
	}
	
	// Update Prometheus metrics for queue depth and consumer count
	metrics.UpdateQueueMetrics(tenantID, queueName, float64(queue.Messages), float64(queue.Consumers))

	// Set default values
	status := "inactive"
	workers := 0

	// Update with actual values if consumer exists
	if consumer != nil {
		status = "active"
		workers = int(consumer.WorkerCount.Load())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":         status,
		"workers":        workers,
		"messageCount":   queue.Messages,
		"consumerCount":  queue.Consumers,
		"processingRate": "N/A", // TODO: Implement processing rate calculation
	})
}

// PublishMessage handles publishing a message to RabbitMQ for a tenant
func (h *TenantHandler) PublishMessage(c echo.Context) error {
	tenantID := c.Param("id")
	if tenantID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "tenant ID is required"})
	}

	// Parse request body
	var message map[string]interface{}
	if err := c.Bind(&message); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid message format"})
	}

	// Get channel from RabbitMQ
	ch, err := h.tenantUseCase.GetChannel()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get channel"})
	}
	defer ch.Close()

	// Convert message to JSON
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to marshal message"})
	}

	// Publish message to RabbitMQ
	exchange := "" // Use default exchange
	routingKey := fmt.Sprintf("tenant.%s", tenantID)

	err = ch.Publish(
		exchange,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        messageBytes,
		},
	)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to publish message"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":   "Message published successfully",
		"tenant_id": tenantID,
	})
}