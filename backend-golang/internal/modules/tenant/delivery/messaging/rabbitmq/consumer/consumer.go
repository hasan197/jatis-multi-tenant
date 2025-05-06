package consumer

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/streadway/amqp"
	"sample-stack-golang/internal/modules/tenant/domain"
	"sample-stack-golang/pkg/logger"
	"sample-stack-golang/pkg/rabbitmq"
)

// StartConsumer memulai consumer untuk tenant tertentu
func StartConsumer(
	ctx context.Context,
	tenantID string,
	rabbitConn *amqp.Connection,
	db *pgxpool.Pool,
	addToWaitGroup func(),
	startWorkerFunc func(*domain.TenantConsumer, int),
) (*domain.TenantConsumer, error) {
	// Buat channel
	ch, err := rabbitConn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %v", err)
	}

	// Setup dead letter configuration
	dlConfig := rabbitmq.NewDefaultDeadLetterConfig()
	
	// Setup dead letter exchange
	err = rabbitmq.SetupDeadLetterExchange(ch, dlConfig)
	if err != nil {
		ch.Close()
		return nil, err
	}
	
	// Setup dead letter queue for tenant
	routingKey, err := rabbitmq.SetupDeadLetterQueue(ch, tenantID, dlConfig)
	if err != nil {
		ch.Close()
		return nil, err
	}
	
	// Declare main queue with dead-letter configuration
	queueName := fmt.Sprintf("tenant.%s", tenantID)
	args := rabbitmq.GetDeadLetterArgs(dlConfig.ExchangeName, routingKey, dlConfig.MessageTTL)
	q, err := ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		args,  // arguments with dead-letter configuration
	)
	if err != nil {
		ch.Close()
		return nil, fmt.Errorf("failed to declare queue: %v", err)
	}

	// Get tenant details from database to determine worker count
	var workerCount int
	if err := db.QueryRow(ctx, "SELECT workers FROM tenants WHERE id = $1", tenantID).Scan(&workerCount); err != nil {
		logger.Log.WithFields(map[string]interface{}{
			"tenant_id": tenantID,
			"error":     err,
		}).Warn("Failed to get worker count from database, using default")
		workerCount = 3 // Default worker count
	}

	// Ensure worker count is at least 1
	if workerCount < 1 {
		workerCount = 1
	}

	// Create buffered message channel for worker pool
	messageChan := make(chan amqp.Delivery, workerCount*10) // Buffer size is 10x worker count

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
		MessageChan:   messageChan,
	}

	// Initialize worker count atomic variable
	consumer.WorkerCount.Store(int32(workerCount))

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
		return nil, fmt.Errorf("failed to start consuming: %v", err)
	}

	// Start message forwarding goroutine
	go forwardMessages(consumer, msgs)

	// Start worker pool
	for i := 0; i < workerCount; i++ {
		workerID := i
		// Add to waitgroup if provided
		if addToWaitGroup != nil {
			addToWaitGroup()
		}
		go startWorkerFunc(consumer, workerID)
	}

	logger.Log.WithFields(map[string]interface{}{
		"tenant_id":    tenantID,
		"worker_count": workerCount,
	}).Info("Started consumer with worker pool")

	return consumer, nil
}

// forwardMessages meneruskan pesan dari RabbitMQ ke message channel untuk diproses oleh worker pool
func forwardMessages(consumer *domain.TenantConsumer, msgs <-chan amqp.Delivery) {
	for {
		select {
		case <-consumer.StopChannel:
			logger.Log.WithFields(map[string]interface{}{
				"tenant_id": consumer.TenantID,
			}).Info("Stopping message forwarding")
			return
		case msg, ok := <-msgs:
			if !ok {
				// Channel closed by RabbitMQ
				logger.Log.WithFields(map[string]interface{}{
					"tenant_id": consumer.TenantID,
				}).Warn("RabbitMQ delivery channel closed unexpectedly")
				return
			}

			// Forward message to worker pool
			select {
			case consumer.MessageChan <- msg:
				// Message forwarded successfully
			case <-consumer.StopChannel:
				// Consumer is stopping, exit
				logger.Log.WithFields(map[string]interface{}{
					"tenant_id": consumer.TenantID,
				}).Info("Stopping message forwarding during message processing")
				return
			}
		}
	}
}
