package rabbitmq

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/streadway/amqp"
	"sample-stack-golang/internal/modules/tenant/domain"
	"sample-stack-golang/pkg/logger"
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
			logger.Log.WithFields(map[string]interface{}{
				"tenant_id": id,
				"error": err,
			}).Error("Error stopping consumer")
		}
	}

	return nil
}

// StartConsumer memulai consumer untuk tenant tertentu
func (m *TenantManager) StartConsumer(ctx context.Context, tenantID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check jika consumer sudah ada
	// if _, exists := m.consumers[tenantID]; exists {
	// 	return fmt.Errorf("consumer already exists for tenant %s", tenantID)
	// }

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

// DebugRabbitMQState prints detailed information about RabbitMQ state
func (m *TenantManager) DebugRabbitMQState(ctx context.Context, tenantID string) {
	logger.Log.WithField("tenant_id", tenantID).Info("=== DEBUG: RabbitMQ State ===")
	
	// Log internal state
	m.mu.RLock()
	consumer, exists := m.consumers[tenantID]
	allConsumers := make(map[string]bool)
	for id := range m.consumers {
		allConsumers[id] = true
	}
	m.mu.RUnlock()
	
	logger.Log.WithFields(map[string]interface{}{
		"tenant_id": tenantID,
		"consumer_exists_in_map": exists,
		"all_consumers": allConsumers,
	}).Info("TenantManager internal state")
	
	if exists {
		logger.Log.WithFields(map[string]interface{}{
			"tenant_id": tenantID,
			"consumer_tag": consumer.ConsumerTag,
			"is_active": consumer.IsActive,
			"channel_nil": consumer.Channel == nil,
		}).Info("Consumer details")
	}
	
	// Try to get RabbitMQ state
	ch, err := m.rabbitConn.Channel()
	if err != nil {
		logger.Log.WithFields(map[string]interface{}{
			"tenant_id": tenantID,
			"error": err.Error(),
		}).Error("Failed to open channel for debug")
		return
	}
	defer ch.Close()
	
	// Check if queue exists
	queueName := fmt.Sprintf("tenant.%s", tenantID)
	queue, err := ch.QueueInspect(queueName)
	if err != nil {
		logger.Log.WithFields(map[string]interface{}{
			"tenant_id": tenantID,
			"queue_name": queueName,
			"error": err.Error(),
		}).Info("Queue does not exist or cannot be inspected")
	} else {
		logger.Log.WithFields(map[string]interface{}{
			"tenant_id": tenantID,
			"queue_name": queue.Name,
			"queue_messages": queue.Messages,
			"queue_consumers": queue.Consumers,
		}).Info("Queue exists")
	}
	
	logger.Log.Info("=== END DEBUG: RabbitMQ State ===")
}

// stopConsumer menghentikan consumer (internal method, assumes lock is held)
func (m *TenantManager) stopConsumer(ctx context.Context, tenantID string) error {
	logger.Log.WithField("tenant_id", tenantID).Info("Starting consumer stop process")
	
	// Cek dan dapatkan consumer
	consumer, err := m.getAndValidateConsumer(tenantID)
	if err != nil {
		return err
	}

	// Stop consumer dan channel
	if err := m.stopConsumerAndChannel(tenantID, consumer); err != nil {
		logger.Log.WithError(err).Warn("Error stopping consumer and channel")
	}

	// Delete queue
	if err := m.deleteQueue(tenantID); err != nil {
		return err
	}

	// Hapus consumer dari map
	m.removeConsumerFromMap(tenantID)

	// Verifikasi queue sudah terhapus
	m.verifyQueueDeletion(tenantID)

	return nil
}

// getAndValidateConsumer mendapatkan dan memvalidasi consumer
func (m *TenantManager) getAndValidateConsumer(tenantID string) (*domain.TenantConsumer, error) {
	consumer, exists := m.consumers[tenantID]
	if !exists {
		logger.Log.WithField("tenant_id", tenantID).Warn("Consumer not found in internal map")
		return nil, fmt.Errorf("consumer not found for tenant %s", tenantID)
	}
	
	logger.Log.WithFields(map[string]interface{}{
		"tenant_id": tenantID,
		"consumer_tag": consumer.ConsumerTag,
		"is_active": consumer.IsActive,
	}).Info("Found consumer in internal map")

	return consumer, nil
}

// stopConsumerAndChannel menghentikan consumer dan menutup channel
func (m *TenantManager) stopConsumerAndChannel(tenantID string, consumer *domain.TenantConsumer) error {
	// Signal stop ke message processing goroutine
	logger.Log.WithField("tenant_id", tenantID).Info("Signaling stop to consumer goroutine")
	close(consumer.StopChannel)

	if consumer.Channel == nil {
		return nil
	}

	// Cancel consumer
	logger.Log.WithField("tenant_id", tenantID).Info("Attempting to cancel consumer")
	if err := consumer.Channel.Cancel(consumer.ConsumerTag, false); err != nil {
		logger.Log.WithFields(map[string]interface{}{
			"tenant_id": tenantID,
			"error":     err,
		}).Warn("Failed to cancel consumer, will force close channel")
	} else {
		logger.Log.WithField("tenant_id", tenantID).Info("Successfully canceled consumer")
	}

	// Close channel
	logger.Log.WithField("tenant_id", tenantID).Info("Closing channel")
	if err := consumer.Channel.Close(); err != nil {
		logger.Log.WithFields(map[string]interface{}{
			"tenant_id": tenantID,
			"error":     err,
		}).Warn("Failed to close channel, continuing with cleanup")
		return err
	}
	
	logger.Log.WithField("tenant_id", tenantID).Info("Successfully closed channel")
	return nil
}

// deleteQueue menghapus queue dari RabbitMQ
func (m *TenantManager) deleteQueue(tenantID string) error {
	logger.Log.WithField("tenant_id", tenantID).Info("Creating channel for queue deletion")
	
	ch, err := m.rabbitConn.Channel()
	if err != nil {
		logger.Log.WithFields(map[string]interface{}{
			"tenant_id": tenantID,
			"error":     err,
		}).Error("Failed to open channel for queue deletion")
		return fmt.Errorf("failed to open channel for queue deletion: %v", err)
	}
	defer ch.Close()

	queueName := fmt.Sprintf("tenant.%s", tenantID)
	logger.Log.WithFields(map[string]interface{}{
		"tenant_id":  tenantID,
		"queue_name": queueName,
	}).Info("Attempting to delete queue")
	
	_, err = ch.QueueDelete(
		queueName, // queue name
		false,     // ifUnused - delete even if queue has consumers
		false,     // ifEmpty - delete even if queue has messages
		false,     // noWait - wait for server response
	)
	
	if err != nil {
		logger.Log.WithFields(map[string]interface{}{
			"tenant_id":  tenantID,
			"queue_name": queueName,
			"error":      err,
		}).Error("Failed to delete queue")
		return fmt.Errorf("failed to delete queue: %v", err)
	}
	
	logger.Log.WithFields(map[string]interface{}{
		"tenant_id":  tenantID,
		"queue_name": queueName,
	}).Info("Successfully deleted queue")

	return nil
}

// removeConsumerFromMap menghapus consumer dari internal map
func (m *TenantManager) removeConsumerFromMap(tenantID string) {
	logger.Log.WithField("tenant_id", tenantID).Info("Removing consumer from internal map")
	delete(m.consumers, tenantID)
}

// verifyQueueDeletion memverifikasi bahwa queue sudah terhapus
func (m *TenantManager) verifyQueueDeletion(tenantID string) {
	defer func() {
		if r := recover(); r != nil {
			logger.Log.WithFields(map[string]interface{}{
				"tenant_id": tenantID,
				"panic":     r,
			}).Warn("Panic during queue verification")
		}
	}()
	
	verifyChannel, err := m.rabbitConn.Channel()
	if err != nil {
		logger.Log.WithFields(map[string]interface{}{
			"tenant_id": tenantID,
			"error":     err,
		}).Warn("Could not verify queue deletion - failed to open channel")
		return
	}
	defer verifyChannel.Close()

	queueName := fmt.Sprintf("tenant.%s", tenantID)
	_, err = verifyChannel.QueueInspect(queueName)
	if err != nil {
		logger.Log.WithFields(map[string]interface{}{
			"tenant_id":  tenantID,
			"queue_name": queueName,
		}).Info("Verified queue has been successfully removed")
	} else {
		logger.Log.WithFields(map[string]interface{}{
			"tenant_id":  tenantID,
			"queue_name": queueName,
		}).Warn("Queue still exists after deletion attempt!")
	}
}

// GetConsumer mendapatkan consumer untuk tenant tertentu
func (m *TenantManager) GetConsumer(tenantID string) *domain.TenantConsumer {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.consumers[tenantID]
}

// logRabbitMQState logs the current state of RabbitMQ connections, channels, and consumers
func (m *TenantManager) logRabbitMQState(ctx context.Context, tenantID string) {
	// Log internal state first
	m.mu.RLock()
	consumer, exists := m.consumers[tenantID]
	m.mu.RUnlock()

	logFields := map[string]interface{}{
		"tenant_id": tenantID,
		"consumer_exists_in_map": exists,
	}

	if exists {
		logFields["consumer_tag"] = consumer.ConsumerTag
		logFields["is_active"] = consumer.IsActive
		logFields["channel_nil"] = consumer.Channel == nil
	}

	// Try to get RabbitMQ state
	ch, err := m.rabbitConn.Channel()
	if err != nil {
		logFields["error_getting_channel"] = err.Error()
		logger.Log.WithFields(logFields).Info("RabbitMQ state - failed to get channel")
		return
	}
	defer ch.Close()

	// Check if queue exists
	queueName := fmt.Sprintf("tenant.%s", tenantID)
	queue, err := ch.QueueInspect(queueName)
	if err != nil {
		logFields["queue_error"] = err.Error()
		logFields["queue_exists"] = false
	} else {
		logFields["queue_exists"] = true
		logFields["queue_name"] = queue.Name
		logFields["queue_messages"] = queue.Messages
		logFields["queue_consumers"] = queue.Consumers
	}

	logger.Log.WithFields(logFields).Info("RabbitMQ state")
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
			logger.Log.WithFields(map[string]interface{}{
				"tenant_id": consumer.TenantID,
				"message": string(msg.Body),
			}).Info("Processing message for tenant")
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
				logger.Log.WithError(err).Error("RabbitMQ health check failed")
				continue
			}
			ch.Close()

			// Check semua consumers
			m.mu.RLock()
			for id, consumer := range m.consumers {
				if !consumer.IsActive {
					logger.Log.WithField("tenant_id", id).Warn("Consumer inactive for tenant, attempting to restart")
					m.mu.RUnlock()
					if err = m.restartConsumer(ctx, id); err != nil {
						logger.Log.WithFields(map[string]interface{}{
							"tenant_id": id,
							"error": err,
						}).Error("Failed to restart consumer for tenant")
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