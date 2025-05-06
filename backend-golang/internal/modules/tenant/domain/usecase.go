package domain

import (
	"context"
)

// TenantManager interface untuk mengelola tenant consumers
type TenantManager interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	StartConsumer(ctx context.Context, tenantID string) error
	StopConsumer(ctx context.Context, tenantID string) error
	GetConsumer(tenantID string) *TenantConsumer
	GetAllConsumers() []*TenantConsumer
	GetActiveConsumers() map[string]*TenantConsumer
	AddConsumer(tenantID string, consumer *TenantConsumer)
	RemoveConsumer(tenantID string)
	UpdateHeartbeat(tenantID string)
	DebugRabbitMQState(ctx context.Context, tenantID string)
}

// TenantUseCase interface untuk business logic tenant
type TenantUseCase interface {
	Create(ctx context.Context, tenant *Tenant) error
	GetByID(ctx context.Context, id string) (*Tenant, error)
	Update(ctx context.Context, tenant *Tenant) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*Tenant, error)
	StartConsumer(ctx context.Context, tenantID string) error
	StopConsumer(ctx context.Context, tenantID string) error
	GetConsumers(ctx context.Context) ([]*TenantConsumer, error)
	UpdateConcurrency(ctx context.Context, id string, config *ConcurrencyConfig) error
}
