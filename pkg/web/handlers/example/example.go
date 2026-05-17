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

// GetExamples returns a paginated list of examples.
//
//	@Summary	List examples
//	@Tags		Example
//	@Produce	json
//	@Param		skip	query		int	false	"Number of records to skip"		default(0)
//	@Param		take	query		int	false	"Number of records to return"	default(20)	maximum(100)
//	@Success	200		{object}	dtos.ExampleListResponse
//	@Failure	400		{object}	object	"invalid pagination params"
//	@Failure	500		{object}	object	"internal server error"
//	@Router		/example [get]
func (h *Handler) GetExamples() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		skip, take, ok := httputil.PaginationParams(w, r, 20, 100)
		if !ok {
			return
		}

		var total int64
		if err := boot.DB.Model(&models.Example{}).Count(&total).Error; err != nil {
			httputil.InternalError(w, r, err)
			return
		}

		var examples []models.Example
		if err := boot.DB.Offset(skip).Limit(take).Find(&examples).Error; err != nil {
			httputil.InternalError(w, r, err)
			return
		}

		items := make([]dtos.ExampleResponse, len(examples))
		for i, e := range examples {
			items[i] = dtos.ExampleResponse{ID: e.ID, Name: e.Name}
		}

		result := dtos.PaginationResult[dtos.ExampleResponse]{
			Data: items,
			Metadata: dtos.PaginationMeta{
				Skip:    skip,
				Take:    take,
				Total:   int(total),
				HasMore: skip+take < int(total),
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result) //nolint:errcheck
	}
}

// CreateExample creates a new example.
//
//	@Summary	Create example
//	@Tags		Example
//	@Accept		json
//	@Produce	json
//	@Param		request	body		dtos.CreateExampleRequest	true	"Example data"
//	@Success	201		{object}	dtos.ExampleResponse
//	@Failure	400		{object}	object	"invalid request body"
//	@Failure	409		{object}	object	"already exists"
//	@Failure	422		{object}	object	"validation error"
//	@Failure	500		{object}	object	"internal server error"
//	@Router		/example [post]
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
		json.NewEncoder(w).Encode(dtos.ExampleResponse{ID: example.ID, Name: example.Name}) //nolint:errcheck
	}
}
