package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/streadway/amqp"
	"sample-stack-golang/internal/modules/tenant/domain"
)

// TenantManager mengimplementasikan domain.TenantManager untuk RabbitMQ
type TenantManager struct {
	rabbitConn *amqp.Connection
	consumers  map[string]*domain.TenantConsumer
	mu         sync.RWMutex
	stopChan   chan struct{}
}

// NewTenantManager membuat instance baru dari TenantManager
func NewTenantManager(rabbitConn *amqp.Connection) domain.TenantManager {
	return &TenantManager{
		rabbitConn: rabbitConn,
		consumers:  make(map[string]*domain.TenantConsumer),
		stopChan:   make(chan struct{}),
	}
}

// Start memulai tenant manager
func (m *TenantManager) Start(ctx context.Context) error {
	// Start health check goroutine
	go m.healthCheck(ctx)
	return nil
}

// Stop menghentikan tenant manager dan semua consumers
func (m *TenantManager) Stop(ctx context.Context) error {
	// Signal stop ke semua goroutines
	close(m.stopChan)

	// Stop semua consumers
	m.mu.Lock()
	defer m.mu.Unlock()

	for id := range m.consumers {
		if err := m.stopConsumer(ctx, id); err != nil {
			log.Printf("Error stopping consumer %s: %v", id, err)
		}
	}

	return nil
}

// StartConsumer memulai consumer untuk tenant tertentu
func (m *TenantManager) StartConsumer(ctx context.Context, tenantID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check jika consumer sudah ada
	if _, exists := m.consumers[tenantID]; exists {
		return fmt.Errorf("consumer already exists for tenant %s", tenantID)
	}

	// Buat channel
	ch, err := m.rabbitConn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %v", err)
	}

	// Declare queue
	queueName := fmt.Sprintf("tenant.%s", tenantID)
	q, err := ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		ch.Close()
		return fmt.Errorf("failed to declare queue: %v", err)
	}

	// Buat consumer
	consumer := &domain.TenantConsumer{
		TenantID:      tenantID,
		QueueName:     queueName,
		ConsumerTag:   fmt.Sprintf("consumer.%s", tenantID),
		Channel:       ch,
		StopChannel:   make(chan struct{}),
		IsActive:      true,
		LastHeartbeat: time.Now(),
		ErrorChannel:  make(chan error, 1),
	}

	// Start consuming
	msgs, err := ch.Consume(
		q.Name,
		consumer.ConsumerTag,
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		ch.Close()
		return fmt.Errorf("failed to start consuming: %v", err)
	}

	// Simpan consumer
	m.consumers[tenantID] = consumer

	// Start message processing goroutine
	go m.processMessages(consumer, msgs)

	return nil
}

// StopConsumer menghentikan consumer untuk tenant tertentu
func (m *TenantManager) StopConsumer(ctx context.Context, tenantID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.stopConsumer(ctx, tenantID)
}

// stopConsumer menghentikan consumer (internal method, assumes lock is held)
func (m *TenantManager) stopConsumer(ctx context.Context, tenantID string) error {
	consumer, exists := m.consumers[tenantID]
	if !exists {
		return fmt.Errorf("consumer not found for tenant %s", tenantID)
	}

	// Signal stop ke message processing goroutine
	close(consumer.StopChannel)

	// Close channel
	if consumer.Channel != nil {
		if err := consumer.Channel.Close(); err != nil {
			return fmt.Errorf("failed to close channel: %v", err)
		}
	}

	// Hapus consumer
	delete(m.consumers, tenantID)

	return nil
}

// GetConsumer mendapatkan consumer untuk tenant tertentu
func (m *TenantManager) GetConsumer(tenantID string) *domain.TenantConsumer {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.consumers[tenantID]
}

// GetAllConsumers mendapatkan semua active consumers
func (m *TenantManager) GetAllConsumers() []*domain.TenantConsumer {
	m.mu.RLock()
	defer m.mu.RUnlock()

	consumers := make([]*domain.TenantConsumer, 0, len(m.consumers))
	for _, consumer := range m.consumers {
		consumers = append(consumers, consumer)
	}
	return consumers
}

// GetActiveConsumers mendapatkan semua active consumers
func (m *TenantManager) GetActiveConsumers() map[string]*domain.TenantConsumer {
	m.mu.RLock()
	defer m.mu.RUnlock()

	activeConsumers := make(map[string]*domain.TenantConsumer)
	for id, consumer := range m.consumers {
		if consumer.IsActive {
			activeConsumers[id] = consumer
		}
	}
	return activeConsumers
}

// AddConsumer menambahkan consumer ke manager
func (m *TenantManager) AddConsumer(tenantID string, consumer *domain.TenantConsumer) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.consumers[tenantID] = consumer
}

// RemoveConsumer menghapus consumer dari manager
func (m *TenantManager) RemoveConsumer(tenantID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.consumers, tenantID)
}

// UpdateHeartbeat memperbarui last heartbeat time untuk consumer
func (m *TenantManager) UpdateHeartbeat(tenantID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if consumer, exists := m.consumers[tenantID]; exists {
		consumer.LastHeartbeat = time.Now()
	}
}

// processMessages memproses pesan untuk consumer
func (m *TenantManager) processMessages(consumer *domain.TenantConsumer, msgs <-chan amqp.Delivery) {
	for {
		select {
		case <-consumer.StopChannel:
			return
		case msg, ok := <-msgs:
			if !ok {
				return
			}
			// Process message
			log.Printf("Processing message for tenant %s: %s", consumer.TenantID, string(msg.Body))
			msg.Ack(false)
		}
	}
}

// healthCheck melakukan health check secara periodik
func (m *TenantManager) healthCheck(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-m.stopChan:
			return
		case <-ticker.C:
			// Check RabbitMQ connection
			ch, err := m.rabbitConn.Channel()
			if err != nil {
				log.Printf("RabbitMQ health check failed: %v", err)
				continue
			}
			ch.Close()

			// Check semua consumers
			m.mu.RLock()
			for id, consumer := range m.consumers {
				if !consumer.IsActive {
					log.Printf("Consumer inactive for tenant %s, attempting to restart", id)
					m.mu.RUnlock()
					if err := m.restartConsumer(ctx, id); err != nil {
						log.Printf("Failed to restart consumer for tenant %s: %v", id, err)
					}
					m.mu.RLock()
				}
			}
			m.mu.RUnlock()
		}
	}
}

// restartConsumer me-restart consumer
func (m *TenantManager) restartConsumer(ctx context.Context, tenantID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Stop existing consumer
	if err := m.stopConsumer(ctx, tenantID); err != nil {
		return err
	}

	// Start new consumer
	return m.StartConsumer(ctx, tenantID)
} 