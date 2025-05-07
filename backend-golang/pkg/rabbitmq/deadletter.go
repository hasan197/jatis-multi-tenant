package rabbitmq

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
	"github.com/jatis/sample-stack-golang/pkg/logger"
)

const (
	// DefaultMaxRetries adalah jumlah maksimal percobaan untuk memproses pesan
	DefaultMaxRetries = 3

	// DefaultMessageTTL adalah waktu hidup default untuk pesan dalam milidetik (24 jam)
	DefaultMessageTTL = int32(1000 * 60 * 60 * 24)
)

// DeadLetterConfig berisi konfigurasi untuk dead letter queue
type DeadLetterConfig struct {
	// ExchangeName adalah nama dari dead letter exchange
	ExchangeName string

	// QueuePrefix adalah prefix untuk nama dead letter queue
	QueuePrefix string

	// MessageTTL adalah waktu hidup pesan dalam milidetik
	MessageTTL int32

	// MaxRetries adalah jumlah maksimal percobaan untuk memproses pesan
	MaxRetries int32
}

// NewDefaultDeadLetterConfig membuat DeadLetterConfig dengan nilai default
func NewDefaultDeadLetterConfig() *DeadLetterConfig {
	return &DeadLetterConfig{
		ExchangeName: "dlx.tenant",
		QueuePrefix:  "dlq.tenant",
		MessageTTL:   DefaultMessageTTL,
		MaxRetries:   DefaultMaxRetries,
	}
}

// SetupDeadLetterExchange membuat dan mengkonfigurasi dead letter exchange
func SetupDeadLetterExchange(ch *amqp.Channel, config *DeadLetterConfig) error {
	// Declare dead-letter exchange
	err := ch.ExchangeDeclare(
		config.ExchangeName,
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare dead-letter exchange: %w", err)
	}

	return nil
}

// SetupDeadLetterQueue membuat dan mengkonfigurasi dead letter queue untuk tenant tertentu
func SetupDeadLetterQueue(ch *amqp.Channel, tenantID string, config *DeadLetterConfig) (string, error) {
	// Declare dead-letter queue
	dlqName := fmt.Sprintf("%s.%s", config.QueuePrefix, tenantID)
	_, err := ch.QueueDeclare(
		dlqName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return "", fmt.Errorf("failed to declare dead-letter queue: %w", err)
	}

	// Routing key untuk binding
	routingKey := fmt.Sprintf("tenant.%s", tenantID)

	// Bind dead-letter queue to exchange
	err = ch.QueueBind(
		dlqName,
		routingKey,
		config.ExchangeName,
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return "", fmt.Errorf("failed to bind dead-letter queue: %w", err)
	}

	return routingKey, nil
}

// GetDeadLetterArgs mengembalikan arguments untuk queue dengan dead letter configuration
func GetDeadLetterArgs(dlxName, routingKey string, ttl int32) amqp.Table {
	return amqp.Table{
		"x-dead-letter-exchange":    dlxName,
		"x-dead-letter-routing-key": routingKey,
		"x-message-ttl":             ttl,
	}
}

// HandleMessageProcessingError menangani error pemrosesan pesan dengan retry logic
func HandleMessageProcessingError(
	msg amqp.Delivery,
	processingError error,
	queueName string,
	tenantID string,
	workerID int,
	maxRetries int32,
) error {
	// Log awal proses penanganan error
	logger.Log.WithFields(map[string]interface{}{
		"tenant_id":  tenantID,
		"worker_id":  workerID,
		"message_id": msg.MessageId,
		"error":      processingError,
	}).Info("[DLQ] Mulai menangani error pemrosesan pesan")
	// Cek apakah pesan sudah memiliki header x-retry-count
	var retryCount int32 = 0
	if msg.Headers != nil {
		logger.Log.WithFields(map[string]interface{}{
			"tenant_id":  tenantID,
			"worker_id":  workerID,
			"message_id": msg.MessageId,
			"headers":    msg.Headers,
		}).Debug("[DLQ] Memeriksa header pesan")

		if retryCountVal, ok := msg.Headers["x-retry-count"]; ok {
			logger.Log.WithFields(map[string]interface{}{
				"tenant_id":     tenantID,
				"worker_id":     workerID,
				"message_id":    msg.MessageId,
				"retry_count_raw": retryCountVal,
				"type":          fmt.Sprintf("%T", retryCountVal),
			}).Debug("[DLQ] Menemukan header x-retry-count")

			if retryCount, ok = retryCountVal.(int32); !ok {
				// Jika tipe tidak sesuai, coba konversi
				if retryCountFloat, ok := retryCountVal.(float64); ok {
					retryCount = int32(retryCountFloat)
					logger.Log.WithFields(map[string]interface{}{
						"tenant_id":     tenantID,
						"worker_id":     workerID,
						"message_id":    msg.MessageId,
						"retry_count_float": retryCountFloat,
						"retry_count_int":   retryCount,
					}).Debug("[DLQ] Mengkonversi retry count dari float64 ke int32")
				}
			}
		}
	}

	// Increment retry count
	retryCount++

	logger.Log.WithFields(map[string]interface{}{
		"tenant_id":   tenantID,
		"worker_id":   workerID,
		"message_id":  msg.MessageId,
		"retry_count": retryCount,
		"max_retries": maxRetries,
	}).Info("[DLQ] Memeriksa batas retry")

	// Cek apakah sudah mencapai batas retry
	if retryCount <= maxRetries {
		// Reject pesan saat ini
		if rejectErr := msg.Reject(false); rejectErr != nil {
			logger.Log.WithFields(map[string]interface{}{
				"tenant_id":   tenantID,
				"worker_id":   workerID,
				"message_id":  msg.MessageId,
				"error":       rejectErr,
				"retry_count": retryCount,
			}).Error("Failed to reject message for retry")
			return rejectErr
		}

		// Kita perlu memodifikasi pendekatan retry karena amqp.Delivery tidak memiliki field Channel
		// Sebagai solusi, kita akan menggunakan pendekatan lain untuk retry
		
		// Log retry attempt
		delay := time.Duration(1<<retryCount) * time.Second // 2, 4, 8 seconds
		logger.Log.WithFields(map[string]interface{}{
			"tenant_id":   tenantID,
			"worker_id":   workerID,
			"message_id":  msg.MessageId,
			"retry_count": retryCount,
			"delay":       delay,
		}).Info("Message will be retried via NACK with requeue=true")
		
		// Untuk retry, kita akan menggunakan NACK dengan requeue=true
		// Ini akan menyebabkan pesan dikembalikan ke queue dan diproses ulang
		// Namun, ini tidak memberikan delay. Untuk implementasi yang lebih baik,
		// kita perlu menerima channel sebagai parameter
		logger.Log.WithFields(map[string]interface{}{
			"tenant_id":   tenantID,
			"worker_id":   workerID,
			"message_id":  msg.MessageId,
			"retry_count": retryCount,
		}).Info("[DLQ] Melakukan NACK dengan requeue=true untuk retry")

		if nackErr := msg.Nack(false, true); nackErr != nil { // multiple=false, requeue=true
			logger.Log.WithFields(map[string]interface{}{
				"tenant_id":   tenantID,
				"worker_id":   workerID,
				"message_id":  msg.MessageId,
				"error":       nackErr,
				"retry_count": retryCount,
			}).Error("[DLQ] Failed to nack message for retry")
			return nackErr
		}

		logger.Log.WithFields(map[string]interface{}{
			"tenant_id":   tenantID,
			"worker_id":   workerID,
			"message_id":  msg.MessageId,
			"retry_count": retryCount,
		}).Info("[DLQ] Pesan berhasil di-NACK untuk retry")
	} else {
		// Sudah mencapai batas retry, kirim ke dead-letter queue
		logger.Log.WithFields(map[string]interface{}{
			"tenant_id":   tenantID,
			"worker_id":   workerID,
			"message_id":  msg.MessageId,
			"error":       processingError,
			"retry_count": retryCount,
			"max_retries": maxRetries,
		}).Error("[DLQ] Message processing failed after max retries, sending to dead-letter queue")

		// Reject tanpa requeue akan mengirim ke dead-letter queue
		logger.Log.WithFields(map[string]interface{}{
			"tenant_id":   tenantID,
			"worker_id":   workerID,
			"message_id":  msg.MessageId,
			"retry_count": retryCount,
		}).Info("[DLQ] Melakukan REJECT tanpa requeue untuk mengirim ke DLQ")

		if err := msg.Reject(false); err != nil {
			logger.Log.WithFields(map[string]interface{}{
				"tenant_id":  tenantID,
				"worker_id":  workerID,
				"message_id": msg.MessageId,
				"error":      err,
			}).Error("[DLQ] Failed to reject message to dead-letter queue")
			return err
		}

		logger.Log.WithFields(map[string]interface{}{
			"tenant_id":   tenantID,
			"worker_id":   workerID,
			"message_id":  msg.MessageId,
			"retry_count": retryCount,
		}).Info("[DLQ] Pesan berhasil dikirim ke dead-letter queue")
	}

	return nil
}
