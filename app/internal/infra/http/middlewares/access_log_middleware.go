package middlewares

import (
	"context"
	"github.com/google/uuid"
	"net/http"
	"time"
	"userservice/internal/pkg/commons/logger"
)

// responseWriter wraps ResponseWriter in order to capture the status code of the response
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code before calling the underlying ResponseWriter
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// LoggingMiddleware logs request details, response status, and the duration of the request's process
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := generateRequestID()

		log := logger.NewLogger()
		log.Info(context.Background()).Attributes("Incoming request", map[string]interface{}{
			"request_id":  requestID,
			"method":      r.Method,
			"path":        r.URL.Path,
			"remote_addr": r.RemoteAddr,
		})

		start := time.Now()

		wrappedWriter := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // we default to 200
		}
		next.ServeHTTP(wrappedWriter, r)

		duration := time.Since(start)

		log.Info(context.Background()).Attributes("Request completed", map[string]interface{}{
			"request_id":  requestID,
			"method":      r.Method,
			"path":        r.URL.Path,
			"remote_addr": r.RemoteAddr,
			"duration":    duration,
			"status_code": wrappedWriter.statusCode,
		})
	})
}

// generateRequestID provides a unique identifier for each request
func generateRequestID() string {
	return uuid.New().String()
}
