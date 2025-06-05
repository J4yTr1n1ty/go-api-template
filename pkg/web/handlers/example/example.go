package example

import (
	"encoding/json"
	"net/http"

	"github.com/J4yTr1n1ty/go-api-template/pkg/boot"
	"github.com/J4yTr1n1ty/go-api-template/pkg/models"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) GetExamples() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		examples := []models.Example{}
		boot.DB.Find(&examples)

		names := []string{}
		for _, example := range examples {
			names = append(names, example.Name)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(names)
	}
}
