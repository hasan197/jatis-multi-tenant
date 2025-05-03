package usecase

import (
	"context"
	"time"

	"sample-stack-golang/internal/modules/tenant/domain"
)

// TenantUsecase implements tenant business logic
type TenantUsecase struct {
	tenantRepo domain.TenantRepository
}

// NewTenantUsecase creates a new tenant usecase
func NewTenantUsecase(tenantRepo domain.TenantRepository) *TenantUsecase {
	return &TenantUsecase{
		tenantRepo: tenantRepo,
	}
}

// Create creates a new tenant
func (u *TenantUsecase) Create(ctx context.Context, tenant *domain.Tenant) error {
	tenant.CreatedAt = time.Now()
	tenant.UpdatedAt = time.Now()
	return u.tenantRepo.Create(ctx, tenant)
}

// GetByID gets a tenant by ID
func (u *TenantUsecase) GetByID(ctx context.Context, id string) (*domain.Tenant, error) {
	return u.tenantRepo.GetByID(ctx, id)
}

// Update updates a tenant
func (u *TenantUsecase) Update(ctx context.Context, tenant *domain.Tenant) error {
	tenant.UpdatedAt = time.Now()
	return u.tenantRepo.Update(ctx, tenant)
}

// Delete deletes a tenant
func (u *TenantUsecase) Delete(ctx context.Context, id string) error {
	return u.tenantRepo.Delete(ctx, id)
}

// List lists all tenants
func (u *TenantUsecase) List(ctx context.Context) ([]*domain.Tenant, error) {
	return u.tenantRepo.List(ctx)
} 