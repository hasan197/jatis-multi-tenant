package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"sample-stack-golang/internal/modules/tenant/domain"
	"sample-stack-golang/internal/modules/tenant/usecase"
)

// TenantHandler handles HTTP requests for tenant
type TenantHandler struct {
	tenantUsecase *usecase.TenantUsecase
}

// NewTenantHandler creates a new tenant handler
func NewTenantHandler(tenantUsecase *usecase.TenantUsecase) *TenantHandler {
	return &TenantHandler{
		tenantUsecase: tenantUsecase,
	}
}

// Create handles tenant creation
func (h *TenantHandler) Create(c echo.Context) error {
	var tenant domain.Tenant
	if err := c.Bind(&tenant); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := h.tenantUsecase.Create(c.Request().Context(), &tenant); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, tenant)
}

// GetByID handles getting a tenant by ID
func (h *TenantHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	tenant, err := h.tenantUsecase.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, tenant)
}

// Update handles tenant update
func (h *TenantHandler) Update(c echo.Context) error {
	id := c.Param("id")
	var tenant domain.Tenant
	if err := c.Bind(&tenant); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	tenant.ID = id
	if err := h.tenantUsecase.Update(c.Request().Context(), &tenant); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, tenant)
}

// Delete handles tenant deletion
func (h *TenantHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.tenantUsecase.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// List handles listing all tenants
func (h *TenantHandler) List(c echo.Context) error {
	tenants, err := h.tenantUsecase.List(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, tenants)
} 