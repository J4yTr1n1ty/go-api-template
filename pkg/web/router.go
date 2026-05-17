package web

import (
	"net/http"

	"github.com/J4yTr1n1ty/go-api-template/pkg/web/handlers/example"
)

func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	exampleHandler := example.NewHandler()
	mux.Handle("GET /example", exampleHandler.GetExamples())

	return mux
}
