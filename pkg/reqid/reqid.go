// Package reqid provides request ID storage and retrieval via context.
package reqid

import "context"

type contextKey struct{}

// WithContext returns a new context carrying the given request ID.
func WithContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, contextKey{}, id)
}

// FromContext retrieves the request ID stored by WithContext.
// Returns an empty string if no ID is present.
func FromContext(ctx context.Context) string {
	id, _ := ctx.Value(contextKey{}).(string)
	return id
}
