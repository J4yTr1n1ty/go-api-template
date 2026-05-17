package dtos

type CreateExampleRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
}
