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