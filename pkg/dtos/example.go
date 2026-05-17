package dtos

type CreateExampleRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
}

type ExampleResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// ExampleListResponse is the swagger-visible concrete type for PaginationResult[ExampleResponse].
type ExampleListResponse struct {
	Data     []ExampleResponse `json:"data"`
	Metadata PaginationMeta    `json:"metadata"`
}
