package rabbitmq

import (
	"context"
	"fmt"
	"time"

	"github.com/jatis/sample-stack-golang/internal/modules/tenant/domain"
	"github.com/jatis/sample-stack-golang/pkg/logger"
)

// getAndValidateConsumer mendapatkan dan memvalidasi consumer
func (m *TenantManager) getAndValidateConsumer(tenantID string) (*domain.TenantConsumer, error) {
	consumer, exists := m.consumers[tenantID]
	if !exists || consumer == nil {
		return nil, fmt.Errorf("consumer not found for tenant %s", tenantID)
	}

	if consumer.Channel == nil {
		return nil, fmt.Errorf("channel is nil for tenant %s", tenantID)
	}

	return consumer, nil
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

	if consumer, exists := m.consumers[tenantID]; exists && consumer != nil {
		consumer.LastHeartbeat = time.Now()
	}
}

// DebugRabbitMQState prints detailed information about RabbitMQ state
func (m *TenantManager) DebugRabbitMQState(ctx context.Context, tenantID string) {
	m.logRabbitMQState(ctx, tenantID)
}

// logRabbitMQState logs the current state of RabbitMQ connections, channels, and consumers
func (m *TenantManager) logRabbitMQState(ctx context.Context, tenantID string) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Log connection state
	if m.rabbitConn == nil {
		logger.Log.Error("RabbitMQ connection is nil")
		return
	}

	// Check if connection is closed
	if m.rabbitConn.IsClosed() {
		logger.Log.Error("RabbitMQ connection is closed")
		return
	}

	// Log consumer state
	if tenantID != "" {
		// Log specific tenant consumer
		consumer, exists := m.consumers[tenantID]
		if !exists || consumer == nil {
			logger.Log.WithFields(map[string]interface{}{
				"tenant_id": tenantID,
			}).Error("Consumer not found")
			return
		}

		// Log consumer details
		logger.Log.WithFields(map[string]interface{}{
			"tenant_id":      tenantID,
			"queue_name":     consumer.QueueName,
			"consumer_tag":   consumer.ConsumerTag,
			"is_active":      consumer.IsActive,
			"last_heartbeat": consumer.LastHeartbeat,
			"worker_count":   consumer.WorkerCount.Load(),
		}).Info("Consumer state")

		// Check channel state
		if consumer.Channel == nil {
			logger.Log.WithFields(map[string]interface{}{
				"tenant_id": tenantID,
			}).Error("Channel is nil")
			return
		}

		// Check if channel is closed
		// Note: There's no direct way to check if a channel is closed in amqp
		// We can only check if operations on the channel return an error
	} else {
		// Log all consumers
		logger.Log.WithFields(map[string]interface{}{
			"consumer_count": len(m.consumers),
		}).Info("All consumers")

		for id, consumer := range m.consumers {
			logger.Log.WithFields(map[string]interface{}{
				"tenant_id":      id,
				"queue_name":     consumer.QueueName,
				"consumer_tag":   consumer.ConsumerTag,
				"is_active":      consumer.IsActive,
				"last_heartbeat": consumer.LastHeartbeat,
				"worker_count":   consumer.WorkerCount.Load(),
			}).Info("Consumer state")
		}
	}
}
