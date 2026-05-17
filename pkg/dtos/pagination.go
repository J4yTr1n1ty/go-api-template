package dtos

// PaginationMeta holds pagination state returned with list responses.
type PaginationMeta struct {
	Skip    int  `json:"skip"`
	Take    int  `json:"take"`
	Total   int  `json:"total"`
	HasMore bool `json:"has_more"`
}

// PaginationResult wraps a paginated list with its metadata.
type PaginationResult[T any] struct {
	Data     []T            `json:"data"`
	Metadata PaginationMeta `json:"metadata"`
}
