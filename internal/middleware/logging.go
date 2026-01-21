package middleware

import (
	"bytes"
	"log/slog"
	"net/http"
	"time"

	"github.com/VishalHilal/e-commerce-api/internal/json"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

func (lrw *loggingResponseWriter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	lrw.body.Write(b)
	return lrw.ResponseWriter.Write(b)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create response writer to capture response
		lrw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Log request
		slog.Info("Request started",
			"method", r.Method,
			"path", r.URL.Path,
			"query", r.URL.RawQuery,
			"remote_addr", r.RemoteAddr,
			"user_agent", r.Header.Get("User-Agent"),
		)

		// Process request
		next.ServeHTTP(lrw, r)

		// Calculate duration
		duration := time.Since(start)

		// Log response
		slog.Info("Request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"status", lrw.statusCode,
			"duration_ms", duration.Milliseconds(),
			"response_size", lrw.body.Len(),
		)

		// Log errors for non-2xx status codes
		if lrw.statusCode >= 400 {
			slog.Error("Request error",
				"method", r.Method,
				"path", r.URL.Path,
				"status", lrw.statusCode,
				"response_body", lrw.body.String(),
			)
		}
	})
}

func ErrorLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("Panic recovered",
					"error", err,
					"method", r.Method,
					"path", r.URL.Path,
					"stack", getStackTrace(),
				)

				json.WriteError(w, http.StatusInternalServerError, "Internal server error")
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func getStackTrace() string {
	// In a real implementation, you'd use a proper stack trace library
	// This is a simplified version
	return "stack trace available"
}

// RequestID middleware for tracing
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		w.Header().Set("X-Request-ID", requestID)

		slog.Info("Request with ID",
			"request_id", requestID,
			"method", r.Method,
			"path", r.URL.Path,
		)

		next.ServeHTTP(w, r)
	})
}

func generateRequestID() string {
	// Simple request ID generation
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().Nanosecond()%len(charset)]
	}
	return string(b)
}
