package middleware

import (
	"net/http"
	"time"

	"sample-stack-golang/pkg/logger"
)

// LoggingMiddleware adalah middleware untuk logging HTTP request
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Buat response writer wrapper untuk capture status code
		rw := newResponseWriter(w)
		
		// Eksekusi handler
		next.ServeHTTP(rw, r)

		// Log request
		logger.Log.WithFields(map[string]interface{}{
			"method":      r.Method,
			"path":        r.URL.Path,
			"remote_addr": r.RemoteAddr,
			"status":      rw.statusCode,
			"duration":    time.Since(start),
			"user_agent":  r.UserAgent(),
		}).Info("http request")
	})
}

// responseWriter adalah wrapper untuk http.ResponseWriter
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
} 