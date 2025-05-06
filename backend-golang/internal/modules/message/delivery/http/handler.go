package http

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"sample-stack-golang/internal/modules/message/domain"
	"sample-stack-golang/internal/modules/message/usecase"
)

// MessageHandler handles HTTP requests for message
type MessageHandler struct {
	messageUsecase *usecase.MessageUsecase
}

// NewMessageHandler creates a new message handler
func NewMessageHandler(messageUsecase *usecase.MessageUsecase) *MessageHandler {
	return &MessageHandler{
		messageUsecase: messageUsecase,
	}
}

// Create handles message creation
// @Summary Create a new message
// @Description Create a new message for a specific tenant
// @Tags messages
// @Accept json
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param message body domain.Message true "Message Information"
// @Success 201 {object} domain.Message
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tenants/{tenant_id}/messages [post]
func (h *MessageHandler) Create(c echo.Context) error {
	tenantID, err := uuid.Parse(c.Param("tenant_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid tenant ID"})
	}

	var message domain.Message
	if err := c.Bind(&message); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	message.TenantID = tenantID

	if err := h.messageUsecase.Create(c.Request().Context(), &message); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, message)
}

// GetByID handles getting a message by ID
// @Summary Get message by ID
// @Description Get a specific message by its ID for a tenant
// @Tags messages
// @Accept json
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param id path string true "Message ID"
// @Success 200 {object} domain.Message
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /tenants/{tenant_id}/messages/{id} [get]
func (h *MessageHandler) GetByID(c echo.Context) error {
	tenantID, err := uuid.Parse(c.Param("tenant_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid tenant ID"})
	}

	messageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid message ID"})
	}

	message, err := h.messageUsecase.GetByID(c.Request().Context(), tenantID, messageID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, message)
}

// GetByTenant handles getting messages by tenant ID
// @Summary Get messages by tenant
// @Description Get all messages for a specific tenant with pagination
// @Tags messages
// @Accept json
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param limit query int false "Number of messages to return (default: 10, max: 100)"
// @Param cursor query string false "Cursor for pagination"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tenants/{tenant_id}/messages [get]
func (h *MessageHandler) GetByTenant(c echo.Context) error {
	tenantID, err := uuid.Parse(c.Param("tenant_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid tenant ID"})
	}

	// Parse limit from query param, default to 10 if not provided or invalid
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit <= 0 || limit > 100 {
		limit = 10 // default limit
	}

	filter := domain.MessageFilter{
		TenantID: tenantID,
		Cursor:   c.QueryParam("cursor"),
		Limit:    limit,
	}

	messages, nextCursor, err := h.messageUsecase.GetByTenant(c.Request().Context(), filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":        messages,
		"next_cursor": nextCursor,
	})
}

// GetMessages handles global message retrieval with cursor pagination
// @Summary Get all messages
// @Description Get messages from all tenants with pagination
// @Tags messages
// @Accept json
// @Produce json
// @Param limit query int false "Number of messages to return (default: 10, max: 100)"
// @Param cursor query string false "Cursor for pagination"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /messages [get]
func (h *MessageHandler) GetMessages(c echo.Context) error {
	// Parse limit from query param, default to 10 if not provided or invalid
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit <= 0 || limit > 100 {
		limit = 10 // default limit
	}

	// Get cursor from query param
	cursor := c.QueryParam("cursor")

	// Get messages from all tenants with pagination
	messages, nextCursor, err := h.messageUsecase.GetMessages(c.Request().Context(), cursor, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Return response in the format specified by the task
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":        messages,
		"next_cursor": nextCursor,
	})
}

// Update handles message update
// @Summary Update a message
// @Description Update an existing message for a tenant
// @Tags messages
// @Accept json
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param id path string true "Message ID"
// @Param message body domain.Message true "Updated Message Information"
// @Success 200 {object} domain.Message
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tenants/{tenant_id}/messages/{id} [put]
func (h *MessageHandler) Update(c echo.Context) error {
	tenantID, err := uuid.Parse(c.Param("tenant_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid tenant ID"})
	}

	messageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid message ID"})
	}

	var message domain.Message
	if err := c.Bind(&message); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	message.ID = messageID
	message.TenantID = tenantID

	if err := h.messageUsecase.Update(c.Request().Context(), &message); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, message)
}

// Delete handles message deletion
// @Summary Delete a message
// @Description Delete a message for a tenant
// @Tags messages
// @Accept json
// @Produce json
// @Param tenant_id path string true "Tenant ID"
// @Param id path string true "Message ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tenants/{tenant_id}/messages/{id} [delete]
func (h *MessageHandler) Delete(c echo.Context) error {
	tenantID, err := uuid.Parse(c.Param("tenant_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid tenant ID"})
	}

	messageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid message ID"})
	}

	if err := h.messageUsecase.Delete(c.Request().Context(), tenantID, messageID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
} 