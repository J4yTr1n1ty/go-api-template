package middleware

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/J4yTr1n1ty/go-api-template/pkg/config"
)

// Recovery recovery middleware.
func Recovery(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("PANIC: %v", err)
				debug.PrintStack()
				if config.IsDebug() {
					stack := string(debug.Stack())
					finalMessage := fmt.Sprintf("Internal Server Error: %s:\n%s", err, stack)
					http.Error(w, finalMessage, http.StatusInternalServerError)
				} else {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}
		}()
		handler.ServeHTTP(w, r)
	})
}
