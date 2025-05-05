package http

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"sample-stack-golang/internal/modules/tenant/domain"
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

	// Stop consumer first
	if err := h.tenantUseCase.StopConsumer(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Delete tenant
	if err := h.tenantUseCase.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

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