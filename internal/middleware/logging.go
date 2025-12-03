package middleware

import (
	"log"
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
		log.Printf(
			"method=%s path=%s status=%d duration_ms=%d ip=%s user_agent=%q bytes=%d db_queries=%d db_selects=%d db_inserts=%d db_updates=%d db_deletes=%d",
			r.Method,
			r.URL.Path,
			wrapped.statusCode,
			duration.Milliseconds(),
			r.RemoteAddr,
			r.UserAgent(),
			wrapped.written,
			total,
			selects,
			inserts,
			updates,
			deletes,
		)
	})
}
