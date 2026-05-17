package web

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/J4yTr1n1ty/go-api-template/pkg/web/handlers/example"
	"github.com/J4yTr1n1ty/go-api-template/pkg/web/middleware"
)

func scope(scopes ...middleware.Scope) middleware.Middleware {
	return middleware.RequireScope(scopes...)
}

func wrap(s middleware.Middleware, h http.HandlerFunc) http.Handler {
	return s(http.HandlerFunc(h))
}

func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	exampleHandler := example.NewHandler()
	mux.Handle("GET /example", exampleHandler.GetExamples())
	mux.Handle("POST /example", exampleHandler.CreateExample())

	return mux
}
