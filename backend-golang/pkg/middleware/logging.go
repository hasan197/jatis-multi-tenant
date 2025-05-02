package middleware

import (
	"net/http"
	"time"

	"sample-stack-golang/pkg/logger"
	"go.uber.org/zap"
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
		logger.Log.Info("http request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("remote_addr", r.RemoteAddr),
			zap.Int("status", rw.statusCode),
			zap.Duration("duration", time.Since(start)),
			zap.String("user_agent", r.UserAgent()),
		)
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