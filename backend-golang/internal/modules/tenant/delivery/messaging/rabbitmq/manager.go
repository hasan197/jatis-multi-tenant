package rabbitmq

import (
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/streadway/amqp"
	"github.com/jatis/sample-stack-golang/internal/modules/tenant/domain"
	"github.com/jatis/sample-stack-golang/pkg/graceful"
)

// TenantManager mengimplementasikan domain.TenantManager untuk RabbitMQ
type TenantManager struct {
	rabbitConn      *amqp.Connection
	consumers       map[string]*domain.TenantConsumer
	mu              sync.RWMutex
	stopChan        chan struct{}
	db              *pgxpool.Pool
	shutdownManager *graceful.ShutdownManager
}

// NewTenantManager membuat instance baru dari TenantManager
func NewTenantManager(rabbitConn *amqp.Connection, db *pgxpool.Pool) domain.TenantManager {
	return &TenantManager{
		rabbitConn: rabbitConn,
		consumers:  make(map[string]*domain.TenantConsumer),
		stopChan:   make(chan struct{}),
		db:         db,
	}
}

// SetShutdownManager sets the shutdown manager for graceful shutdown
func (m *TenantManager) SetShutdownManager(sm *graceful.ShutdownManager) {
	m.shutdownManager = sm
}

// GetChannel gets a new channel from RabbitMQ connection
func (m *TenantManager) GetChannel() (*amqp.Channel, error) {
	return m.rabbitConn.Channel()
}
