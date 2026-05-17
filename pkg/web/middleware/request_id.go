package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/J4yTr1n1ty/go-api-template/pkg/reqid"
)

// RequestID generates or propagates a request ID via X-Request-ID header and context.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = generateRequestID()
		}
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(reqid.WithContext(r.Context(), id)))
	})
}

func generateRequestID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
