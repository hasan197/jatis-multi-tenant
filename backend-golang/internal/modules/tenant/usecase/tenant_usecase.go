package usecase

import (
	"context"
	"errors"
	"fmt"

	"sample-stack-golang/internal/modules/tenant/domain"
	"sample-stack-golang/pkg/logger"
)

var (
	ErrTenantNotFound = errors.New("tenant not found")
	ErrInvalidInput   = errors.New("invalid input")
)

// TenantUseCase implements domain.TenantUseCase
type TenantUseCase struct {
	repo    domain.TenantRepository
	manager domain.TenantManager
}

// NewTenantUseCase creates a new tenant usecase
func NewTenantUseCase(repo domain.TenantRepository, manager domain.TenantManager) domain.TenantUseCase {
	return &TenantUseCase{
		repo:    repo,
		manager: manager,
	}
}

// Create creates a new tenant
func (u *TenantUseCase) Create(ctx context.Context, tenant *domain.Tenant) error {
	if err := u.repo.Create(ctx, tenant); err != nil {
		return fmt.Errorf("failed to create tenant: %v", err)
	}

	// Start consumer untuk tenant baru
	if err := u.manager.StartConsumer(ctx, tenant.ID); err != nil {
		// Log error tapi jangan return error karena tenant sudah dibuat
		fmt.Printf("Warning: failed to start consumer for tenant %s: %v\n", tenant.ID, err)
	}

	return nil
}

// GetByID gets a tenant by ID
func (u *TenantUseCase) GetByID(ctx context.Context, id string) (*domain.Tenant, error) {
	tenant, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %v", err)
	}
	if tenant == nil {
		return nil, ErrTenantNotFound
	}
	return tenant, nil
}

// Update updates a tenant
func (u *TenantUseCase) Update(ctx context.Context, tenant *domain.Tenant) error {
	if err := u.repo.Update(ctx, tenant); err != nil {
		return fmt.Errorf("failed to update tenant: %v", err)
	}
	return nil
}

// Delete deletes a tenant
func (u *TenantUseCase) Delete(ctx context.Context, id string) error {
	// Check if tenant exists before proceeding
	tenant, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get tenant: %v", err)
	}
	if tenant == nil {
		return ErrTenantNotFound
	}

	// First, try to stop and remove the consumer
	logger.Log.WithField("tenant_id", id).Info("Attempting to stop consumer before deleting tenant")
	
	// Call DebugRabbitMQState to get detailed information about the RabbitMQ state
	if u.manager != nil {
		u.manager.DebugRabbitMQState(ctx, id)
	}
	
	// Try to stop the consumer
	if u.manager != nil {
		err := u.manager.StopConsumer(ctx, id)
		if err != nil {
			logger.Log.WithFields(map[string]interface{}{
				"tenant_id": id,
				"error":     err,
			}).Warn("Failed to stop consumer, continuing with tenant deletion")
			// Don't return error here, continue with deletion
		} else {
			logger.Log.WithField("tenant_id", id).Info("Successfully stopped consumer")
		}
	}

	// Delete tenant from repository (this will also drop the message partition)
	if err := u.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete tenant: %v", err)
	}

	// Final check to ensure consumer is removed
	if u.manager != nil {
		u.manager.DebugRabbitMQState(ctx, id)
	}

	return nil
}

// List lists all tenants
func (u *TenantUseCase) List(ctx context.Context) ([]*domain.Tenant, error) {
	tenants, err := u.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list tenants: %v", err)
	}
	return tenants, nil
}

// StartConsumer starts a consumer for a tenant
func (u *TenantUseCase) StartConsumer(ctx context.Context, tenantID string) error {
	if err := u.manager.StartConsumer(ctx, tenantID); err != nil {
		return fmt.Errorf("failed to start consumer: %v", err)
	}
	return nil
}

// StopConsumer stops a consumer for a tenant
func (u *TenantUseCase) StopConsumer(ctx context.Context, tenantID string) error {
	if err := u.manager.StopConsumer(ctx, tenantID); err != nil {
		return fmt.Errorf("failed to stop consumer: %v", err)
	}
	return nil
}

// GetConsumers gets all consumers
func (u *TenantUseCase) GetConsumers(ctx context.Context) ([]*domain.TenantConsumer, error) {
	consumers := u.manager.GetAllConsumers()
	return consumers, nil
}

// StopConsumer stops the consumer
func StopConsumer(consumer *domain.TenantConsumer) error {
	close(consumer.StopChannel)
	return nil
} 