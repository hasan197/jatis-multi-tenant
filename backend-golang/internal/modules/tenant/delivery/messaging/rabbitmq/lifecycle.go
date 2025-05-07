package rabbitmq

import (
	"context"
	"fmt"
	"time"

	"github.com/jatis/sample-stack-golang/internal/modules/tenant/delivery/messaging/rabbitmq/consumer"
	"github.com/jatis/sample-stack-golang/internal/modules/tenant/domain"
	"github.com/jatis/sample-stack-golang/pkg/logger"
)

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

	// Get tenant details to determine worker count
	// First, check if we already have a consumer for this tenant
	oldConsumer, exists := m.consumers[tenantID]
	if exists && oldConsumer != nil {
		// Stop existing consumer first
		if err := m.stopConsumer(ctx, tenantID); err != nil {
			logger.Log.WithFields(map[string]interface{}{
				"tenant_id": tenantID,
				"error":     err,
			}).Warn("Error stopping existing consumer before restart")
		}
	}

	// Start consumer with worker pool
	addToWaitGroup := func() {
		if m.shutdownManager != nil {
			m.shutdownManager.AddTask()
		}
	}

	startWorkerFunc := func(c *domain.TenantConsumer, workerID int) {
		consumer.StartWorker(c, workerID, m.shutdownManager)
	}

	newConsumer, err := consumer.StartConsumer(
		ctx,
		tenantID,
		m.rabbitConn,
		m.db,
		addToWaitGroup,
		startWorkerFunc,
	)

	if err != nil {
		return err
	}

	// Simpan consumer
	m.consumers[tenantID] = newConsumer

	return nil
}

// StopConsumer menghentikan consumer untuk tenant tertentu
func (m *TenantManager) StopConsumer(ctx context.Context, tenantID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.stopConsumer(ctx, tenantID)
}

// stopConsumer menghentikan consumer dan membersihkan resources
func (m *TenantManager) stopConsumer(ctx context.Context, tenantID string) error {
	consumer, err := m.getAndValidateConsumer(tenantID)
	if err != nil {
		return err
	}

	// Stop consumer dan channel
	if err := m.stopConsumerAndChannel(tenantID, consumer); err != nil {
		return err
	}

	// Delete queue
	if err := m.deleteQueue(tenantID); err != nil {
		logger.Log.WithFields(map[string]interface{}{
			"tenant_id": tenantID,
			"error":     err,
		}).Warn("Failed to delete queue")
	}

	// Remove consumer from map
	m.removeConsumerFromMap(tenantID)

	return nil
}

// restartConsumer me-restart consumer
func (m *TenantManager) restartConsumer(ctx context.Context, tenantID string) error {
	// Stop consumer
	if err := m.StopConsumer(ctx, tenantID); err != nil {
		return fmt.Errorf("failed to stop consumer: %w", err)
	}

	// Start consumer
	if err := m.StartConsumer(ctx, tenantID); err != nil {
		return fmt.Errorf("failed to start consumer: %w", err)
	}

	return nil
}

// healthCheck melakukan health check secara periodik
func (m *TenantManager) healthCheck(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-m.stopChan:
			logger.Log.Info("Stopping health check")
			return
		case <-ctx.Done():
			logger.Log.Info("Context cancelled, stopping health check")
			return
		case <-ticker.C:
			m.mu.RLock()
			for id, consumer := range m.consumers {
				// Skip if consumer is not active
				if !consumer.IsActive {
					continue
				}

				// Check if consumer is still active
				if time.Since(consumer.LastHeartbeat) > 60*time.Second {
					logger.Log.WithFields(map[string]interface{}{
						"tenant_id": id,
						"last_heartbeat": consumer.LastHeartbeat,
					}).Warn("Consumer heartbeat timeout, restarting")

					// Restart consumer
					go func(tenantID string) {
						if err := m.restartConsumer(ctx, tenantID); err != nil {
							logger.Log.WithFields(map[string]interface{}{
								"tenant_id": tenantID,
								"error":     err,
							}).Error("Failed to restart consumer")
						}
					}(id)
				}
			}
			m.mu.RUnlock()
		}
	}
}
