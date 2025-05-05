package domain

import (
	"context"
	"time"
)

// Tenant represents a tenant in the system
type Tenant struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Description string    `json:"description"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TenantConsumer represents a RabbitMQ consumer for a tenant
type TenantConsumer struct {
	TenantID      string
	QueueName     string
	ConsumerTag   string
	Channel       interface{}
	StopChannel   chan struct{}
	IsActive      bool
	LastHeartbeat time.Time
	ErrorChannel  chan error
}

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
}

// TenantRepository interface untuk operasi database tenant
type TenantRepository interface {
	Create(ctx context.Context, tenant *Tenant) error
	GetByID(ctx context.Context, id string) (*Tenant, error)
	Update(ctx context.Context, tenant *Tenant) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*Tenant, error)
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
}

// Stop stops the consumer
func (tc *TenantConsumer) Stop() error {
	close(tc.StopChannel)
	return nil
} 