package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type wrappedWriter struct {
	http.ResponseWriter
	status int
}

func (w *wrappedWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &wrappedWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)
		duration := time.Since(start)

		requestGroup := slog.String("method", r.Method)
		urlGroup := slog.String("url", r.URL.Path)
		responseStatus := slog.Int("status", wrapped.status)
		durationGroup := slog.Int64("duration_ms", duration.Milliseconds())

		if wrapped.status >= 400 {
			slog.Error("HTTP request completed with error", requestGroup, urlGroup, responseStatus, durationGroup)
		} else {
			slog.Info("HTTP request completed", requestGroup, urlGroup, responseStatus, durationGroup)
		}
	})
}
