package metrics

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// HTTP metrics
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// Custom business metrics
	activeUsers = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_users",
			Help: "Number of active users",
		},
	)

	// Database metrics
	dbQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Duration of database queries in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"query_type"},
	)

	// RabbitMQ metrics - Queue metrics
	QueueDepth = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rabbitmq_queue_depth",
			Help: "The current number of messages in the queue",
		},
		[]string{"tenant_id", "queue_name"},
	)

	QueueConsumerCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rabbitmq_queue_consumer_count",
			Help: "The current number of consumers for the queue",
		},
		[]string{"tenant_id", "queue_name"},
	)

	// RabbitMQ metrics - Worker metrics
	WorkerCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rabbitmq_worker_count",
			Help: "The current number of active workers",
		},
		[]string{"tenant_id"},
	)

	MessageProcessed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rabbitmq_messages_processed_total",
			Help: "The total number of processed messages",
		},
		[]string{"tenant_id", "status"},
	)

	MessageProcessingTime = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "rabbitmq_message_processing_time_seconds",
			Help:    "Time taken to process messages",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"tenant_id"},
	)

	// RabbitMQ metrics - DLQ metrics
	DLQDepth = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rabbitmq_dlq_depth",
			Help: "The current number of messages in the dead letter queue",
		},
		[]string{"tenant_id"},
	)

	MessageRetryCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rabbitmq_message_retry_total",
			Help: "The total number of message retries",
		},
		[]string{"tenant_id"},
	)

	MessageDeadLettered = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rabbitmq_messages_dead_lettered_total",
			Help: "The total number of messages sent to dead letter queue",
		},
		[]string{"tenant_id"},
	)
)

// SetupMetrics mengatur endpoint metrics dan middleware
func SetupMetrics(e *echo.Echo) {
	// Setup metrics endpoint
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	// Setup middleware untuk mengumpulkan metrics
	e.Use(metricsMiddleware())
}

// metricsMiddleware mengumpulkan metrics untuk setiap request HTTP
func metricsMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := prometheus.NewTimer(httpRequestDuration.WithLabelValues(c.Request().Method, c.Path()))
			
			err := next(c)
			
			status := c.Response().Status
			httpRequestsTotal.WithLabelValues(c.Request().Method, c.Path(), fmt.Sprint(status)).Inc()
			start.ObserveDuration()
			
			return err
		}
	}
}

// RecordDBQueryDuration mencatat durasi query database
func RecordDBQueryDuration(queryType string, duration float64) {
	dbQueryDuration.WithLabelValues(queryType).Observe(duration)
}

// SetActiveUsers mengupdate jumlah user aktif
func SetActiveUsers(count float64) {
	activeUsers.Set(count)
}

// RabbitMQ Metrics Functions

// RecordMessageProcessed increments the counter for processed messages
func RecordMessageProcessed(tenantID, status string) {
	MessageProcessed.WithLabelValues(tenantID, status).Inc()
}

// RecordMessageProcessingTime observes the time taken to process a message
func RecordMessageProcessingTime(tenantID string, durationSeconds float64) {
	MessageProcessingTime.WithLabelValues(tenantID).Observe(durationSeconds)
}

// UpdateQueueMetrics updates the queue depth and consumer count metrics
func UpdateQueueMetrics(tenantID, queueName string, depth, consumerCount float64) {
	QueueDepth.WithLabelValues(tenantID, queueName).Set(depth)
	QueueConsumerCount.WithLabelValues(tenantID, queueName).Set(consumerCount)
}

// UpdateWorkerCount updates the worker count metric
func UpdateWorkerCount(tenantID string, count float64) {
	WorkerCount.WithLabelValues(tenantID).Set(count)
}

// UpdateDLQMetrics updates the DLQ depth metric
func UpdateDLQMetrics(tenantID string, depth float64) {
	DLQDepth.WithLabelValues(tenantID).Set(depth)
}

// RecordMessageRetry increments the counter for message retries
func RecordMessageRetry(tenantID string) {
	MessageRetryCount.WithLabelValues(tenantID).Inc()
}

// RecordMessageDeadLettered increments the counter for dead lettered messages
func RecordMessageDeadLettered(tenantID string) {
	MessageDeadLettered.WithLabelValues(tenantID).Inc()
} 