package example

import (
	"encoding/json"
	"net/http"

	"github.com/J4yTr1n1ty/go-api-template/pkg/boot"
	"github.com/J4yTr1n1ty/go-api-template/pkg/dtos"
	"github.com/J4yTr1n1ty/go-api-template/pkg/models"
	"github.com/J4yTr1n1ty/go-api-template/pkg/validate"
	"github.com/J4yTr1n1ty/go-api-template/pkg/web/httputil"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) GetExamples() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var examples []models.Example
		if err := boot.DB.Find(&examples).Error; err != nil {
			httputil.InternalError(w, r, err)
			return
		}

		names := make([]string, len(examples))
		for i, example := range examples {
			names[i] = example.Name
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(names) //nolint:errcheck
	}
}

func (h *Handler) CreateExample() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dtos.CreateExampleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httputil.ErrorJSON(w, r, http.StatusBadRequest, "invalid request body")
			return
		}

		if err := validate.Struct(req); err != nil {
			httputil.ValidationErrorJSON(w, r, err.(validate.FieldErrors))
			return
		}

		example := models.Example{Name: req.Name}
		if err := boot.DB.Create(&example).Error; err != nil {
			if httputil.DBError(w, r, err) {
				return
			}
			httputil.InternalError(w, r, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(example) //nolint:errcheck
	}
}
