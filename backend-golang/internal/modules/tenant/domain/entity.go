package domain

import (
	"sync/atomic"
	"time"

	"github.com/streadway/amqp"
)

// Tenant represents a tenant in the system
type Tenant struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Workers     int       `json:"workers"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TenantConsumer represents a RabbitMQ consumer for a tenant
type TenantConsumer struct {
	TenantID      string         `json:"tenant_id"`
	QueueName     string         `json:"queue_name"`
	ConsumerTag   string         `json:"consumer_tag"`
	Channel       *amqp.Channel  `json:"-"`
	StopChannel   chan struct{}  `json:"-"`
	IsActive      bool           `json:"is_active"`
	LastHeartbeat time.Time      `json:"last_heartbeat"`
	ErrorChannel  chan error     `json:"-"`
	WorkerCount   atomic.Int32   `json:"worker_count"`
	MessageChan   chan amqp.Delivery `json:"-"`
}

// ConcurrencyConfig represents the concurrency configuration for a tenant
type ConcurrencyConfig struct {
	Workers int `json:"workers"`
}