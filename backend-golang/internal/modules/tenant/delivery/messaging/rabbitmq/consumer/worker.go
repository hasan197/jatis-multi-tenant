package consumer

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/streadway/amqp"
	"sample-stack-golang/internal/modules/tenant/domain"
	"sample-stack-golang/pkg/graceful"
	"sample-stack-golang/pkg/logger"
	"sample-stack-golang/pkg/rabbitmq"
)

// StartWorker memulai worker untuk memproses pesan dari message channel
func StartWorker(consumer *domain.TenantConsumer, workerID int, shutdownManager *graceful.ShutdownManager) {
	// Mark worker as done in waitgroup when finished if shutdown manager is available
	if shutdownManager != nil {
		defer shutdownManager.DoneTask()
	}

	logger.Log.WithFields(map[string]interface{}{
		"tenant_id": consumer.TenantID,
		"worker_id": workerID,
	}).Info("Starting worker")

	for {
		select {
		case <-consumer.StopChannel:
			logger.Log.WithFields(map[string]interface{}{
				"tenant_id": consumer.TenantID,
				"worker_id": workerID,
			}).Info("Worker received stop signal")
			return
		case msg, ok := <-consumer.MessageChan:
			if !ok {
				// Channel closed, exit worker
				logger.Log.WithFields(map[string]interface{}{
					"tenant_id": consumer.TenantID,
					"worker_id": workerID,
				}).Info("Message channel closed, stopping worker")
				return
			}

			// Process message
			logger.Log.WithFields(map[string]interface{}{
				"tenant_id":  consumer.TenantID,
				"worker_id":  workerID,
				"message_id": msg.MessageId,
			}).Debug("Processing message")

			// Simulasi pemrosesan dengan kemungkinan error
			var processingError error
			
			// Saat ini hanya simulasi pemrosesan dengan delay
			time.Sleep(100 * time.Millisecond)

			// Periksa apakah pesan memiliki flag force_error
			var payload map[string]interface{}

			// Dekode body pesan
			if err := json.Unmarshal(msg.Body, &payload); err != nil {
				logger.Log.WithFields(map[string]interface{}{
					"tenant_id":  consumer.TenantID,
					"worker_id":  workerID,
					"message_id": msg.MessageId,
					"error":      err,
				}).Error("Failed to decode message payload")
				
				// Jika gagal decode, reject pesan
				if err := msg.Reject(false); err != nil {
					logger.Log.WithFields(map[string]interface{}{
						"tenant_id":  consumer.TenantID,
						"worker_id":  workerID,
						"message_id": msg.MessageId,
						"error":      err,
					}).Error("Failed to reject message after decode error")
				}
				return
			}

			// Log payload untuk debugging
			logger.Log.WithFields(map[string]interface{}{
				"tenant_id":  consumer.TenantID,
				"worker_id":  workerID,
				"message_id": msg.MessageId,
				"payload":    payload,
			}).Debug("Decoded message payload")

			// Periksa apakah ada metadata.force_error
			if metadata, ok := payload["metadata"].(map[string]interface{}); ok {
				logger.Log.WithFields(map[string]interface{}{
					"tenant_id":  consumer.TenantID,
					"worker_id":  workerID,
					"message_id": msg.MessageId,
					"metadata":   metadata,
				}).Debug("Checking message metadata for force_error flag")

				// Periksa apakah ada flag force_error
				if forceError, ok := metadata["force_error"].(bool); ok && forceError {
					logger.Log.WithFields(map[string]interface{}{
						"tenant_id":  consumer.TenantID,
						"worker_id":  workerID,
						"message_id": msg.MessageId,
					}).Info("Force error flag detected, simulating processing error")
					processingError = fmt.Errorf("forced error for testing DLQ")
				}
			} else {
				logger.Log.WithFields(map[string]interface{}{
					"tenant_id":  consumer.TenantID,
					"worker_id":  workerID,
					"message_id": msg.MessageId,
					"payload":    payload,
				}).Debug("No metadata field found in message payload")
			}

			// Simulasi error acak untuk testing (dalam produksi, ini akan diganti dengan error handling yang sebenarnya)
			// Dalam implementasi nyata, ini akan diganti dengan logika bisnis yang sebenarnya
			// dan error handling yang tepat

			// Jika terjadi error dalam pemrosesan
			if processingError != nil {
				logger.Log.WithFields(map[string]interface{}{
					"tenant_id":  consumer.TenantID,
					"worker_id":  workerID,
					"message_id": msg.MessageId,
					"error":      processingError,
				}).Error("Message processing failed, handling with DLQ mechanism")

				// Gunakan package rabbitmq untuk menangani error pemrosesan pesan
				err := rabbitmq.HandleMessageProcessingError(
					msg,
					processingError,
					consumer.QueueName,
					consumer.TenantID,
					workerID,
					rabbitmq.DefaultMaxRetries,
				)
				
				if err != nil {
					logger.Log.WithFields(map[string]interface{}{
						"tenant_id":  consumer.TenantID,
						"worker_id":  workerID,
						"message_id": msg.MessageId,
						"error":      err,
					}).Error("Failed to handle message processing error")
				} else {
					logger.Log.WithFields(map[string]interface{}{
						"tenant_id":  consumer.TenantID,
						"worker_id":  workerID,
						"message_id": msg.MessageId,
					}).Info("Successfully handled message processing error with DLQ mechanism")
				}
			} else {
				// Jika pemrosesan berhasil, acknowledge message
				logger.Log.WithFields(map[string]interface{}{
					"tenant_id":  consumer.TenantID,
					"worker_id":  workerID,
					"message_id": msg.MessageId,
				}).Debug("Message processed successfully, acknowledging")

				if err := msg.Ack(false); err != nil {
					logger.Log.WithFields(map[string]interface{}{
						"tenant_id":  consumer.TenantID,
						"worker_id":  workerID,
						"message_id": msg.MessageId,
						"error":      err,
					}).Error("Failed to acknowledge message")
				} else {
					logger.Log.WithFields(map[string]interface{}{
						"tenant_id":  consumer.TenantID,
						"worker_id":  workerID,
						"message_id": msg.MessageId,
					}).Debug("Message successfully acknowledged")
				}
			}
		}
	}
}

// ProcessMessage adalah fungsi untuk memproses pesan yang diterima
// Ini adalah template yang dapat diimplementasikan sesuai kebutuhan bisnis
func ProcessMessage(msg amqp.Delivery) error {
	// Implementasi pemrosesan pesan sesuai kebutuhan bisnis
	// Contoh:
	// 1. Parse body pesan (msg.Body)
	// 2. Validasi data
	// 3. Simpan ke database atau proses sesuai kebutuhan bisnis
	
	// Simulasi pemrosesan
	time.Sleep(100 * time.Millisecond)
	
	return nil
}
