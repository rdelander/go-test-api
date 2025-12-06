package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.written += n
	return n, err
}

// Logging middleware logs canonical request information
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := newResponseWriter(w)

		// Add DB stats tracking to context
		ctx := WithDBStats(r.Context())
		r = r.WithContext(ctx)

		// Process request
		next.ServeHTTP(wrapped, r)

		// Get DB stats
		stats := GetDBStats(ctx)
		total, selects, inserts, updates, deletes := stats.Summary()

		// Log canonical line
		duration := time.Since(start)
		slog.Info("http_request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", wrapped.statusCode,
			"duration_ms", duration.Milliseconds(),
			"ip", r.RemoteAddr,
			"user_agent", r.UserAgent(),
			"bytes", wrapped.written,
			"db", slog.GroupValue(
				slog.Int("queries", total),
				slog.Int("selects", selects),
				slog.Int("inserts", inserts),
				slog.Int("updates", updates),
				slog.Int("deletes", deletes),
			),
		)
	})
}
