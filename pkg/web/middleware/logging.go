package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)

		remoteip := r.Header.Get("X-Forwarded-For")
		if remoteip == "" {
			remoteip = r.RemoteAddr
		}

		slog.Info("request", "remote_ip", remoteip, "status", wrapped.statusCode, "method", r.Method, "path", r.URL.Path, "duration", time.Since(start))
	})
}
