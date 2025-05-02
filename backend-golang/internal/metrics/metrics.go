package metrics

import (
	"github.com/gin-gonic/gin"
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
)

// SetupMetrics mengatur endpoint metrics dan middleware
func SetupMetrics(router *gin.Engine) {
	// Setup metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Setup middleware untuk mengumpulkan metrics
	router.Use(metricsMiddleware())
}

// metricsMiddleware mengumpulkan metrics untuk setiap request HTTP
func metricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := prometheus.NewTimer(httpRequestDuration.WithLabelValues(c.Request.Method, c.FullPath()))
		
		c.Next()
		
		status := c.Writer.Status()
		httpRequestsTotal.WithLabelValues(c.Request.Method, c.FullPath(), string(status)).Inc()
		start.ObserveDuration()
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