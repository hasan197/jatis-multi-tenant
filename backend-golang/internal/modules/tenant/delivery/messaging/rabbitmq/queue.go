package rabbitmq

import (
	"fmt"
	"time"

	"sample-stack-golang/internal/modules/tenant/domain"
	"sample-stack-golang/pkg/logger"
)

// stopConsumerAndChannel menghentikan consumer dan menutup channel
func (m *TenantManager) stopConsumerAndChannel(tenantID string, consumer *domain.TenantConsumer) error {
	// Signal stop ke consumer
	close(consumer.StopChannel)

	// Tunggu channel ditutup
	time.Sleep(100 * time.Millisecond)

	// Cancel consumer
	if consumer.Channel != nil {
		if err := consumer.Channel.Cancel(consumer.ConsumerTag, false); err != nil {
			logger.Log.WithFields(map[string]interface{}{
				"tenant_id": tenantID,
				"error":     err,
			}).Warn("Failed to cancel consumer")
		}

		// Close channel
		if err := consumer.Channel.Close(); err != nil {
			logger.Log.WithFields(map[string]interface{}{
				"tenant_id": tenantID,
				"error":     err,
			}).Warn("Failed to close channel")
		}
	}

	return nil
}

// deleteQueue menghapus queue dari RabbitMQ
func (m *TenantManager) deleteQueue(tenantID string) error {
	// Buat channel baru untuk delete queue
	ch, err := m.rabbitConn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel for queue deletion: %w", err)
	}
	defer ch.Close()

	// Delete queue
	queueName := fmt.Sprintf("tenant.%s", tenantID)
	_, err = ch.QueueDelete(
		queueName,
		false, // ifUnused
		false, // ifEmpty
		false, // noWait
	)
	if err != nil {
		return fmt.Errorf("failed to delete queue: %w", err)
	}

	// Delete dead-letter queue
	dlqName := fmt.Sprintf("dlq.tenant.%s", tenantID)
	_, err = ch.QueueDelete(
		dlqName,
		false, // ifUnused
		false, // ifEmpty
		false, // noWait
	)
	if err != nil {
		logger.Log.WithFields(map[string]interface{}{
			"tenant_id": tenantID,
			"error":     err,
		}).Warn("Failed to delete dead-letter queue")
	}

	// Verify queue deletion
	go m.verifyQueueDeletion(tenantID)

	return nil
}

// verifyQueueDeletion memverifikasi bahwa queue sudah terhapus
func (m *TenantManager) verifyQueueDeletion(tenantID string) {
	// Tunggu sebentar untuk memastikan queue sudah terhapus
	time.Sleep(500 * time.Millisecond)

	// Buat channel baru untuk verifikasi
	ch, err := m.rabbitConn.Channel()
	if err != nil {
		logger.Log.WithFields(map[string]interface{}{
			"tenant_id": tenantID,
			"error":     err,
		}).Error("Failed to open channel for queue verification")
		return
	}
	defer ch.Close()

	// Coba declare queue dengan passive=true untuk memeriksa apakah queue masih ada
	queueName := fmt.Sprintf("tenant.%s", tenantID)
	_, err = ch.QueueDeclarePassive(
		queueName,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	)

	if err != nil {
		// Queue tidak ditemukan, ini yang diharapkan
		logger.Log.WithFields(map[string]interface{}{
			"tenant_id": tenantID,
		}).Debug("Queue deletion verified")
	} else {
		// Queue masih ada, ini tidak diharapkan
		logger.Log.WithFields(map[string]interface{}{
			"tenant_id": tenantID,
		}).Warn("Queue still exists after deletion attempt")
	}
}

// removeConsumerFromMap menghapus consumer dari internal map
func (m *TenantManager) removeConsumerFromMap(tenantID string) {
	delete(m.consumers, tenantID)
}
