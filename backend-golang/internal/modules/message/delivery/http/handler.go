package http

import (
	"net/http"

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
func (h *MessageHandler) Create(c echo.Context) error {
	var message domain.Message
	if err := c.Bind(&message); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := h.messageUsecase.Create(c.Request().Context(), &message); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, message)
}

// GetByID handles getting a message by ID
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
func (h *MessageHandler) GetByTenant(c echo.Context) error {
	tenantID, err := uuid.Parse(c.Param("tenant_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid tenant ID"})
	}

	filter := domain.MessageFilter{
		TenantID: tenantID,
		Cursor:   c.QueryParam("cursor"),
		Limit:    10, // default limit
	}

	messages, nextCursor, err := h.messageUsecase.GetByTenant(c.Request().Context(), filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"messages":    messages,
		"next_cursor": nextCursor,
	})
}

// Update handles message update
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