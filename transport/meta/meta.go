package meta

import (
	"context"

	"github.com/alexfalkowski/go-service/meta"
)

const (
	// RequestIDKey for meta.
	RequestIDKey = "requestId"

	// ServiceKey for meta.
	ServiceKey = "service"

	// PathKey for meta.
	PathKey = "path"

	// MethodKey for meta.
	MethodKey = "method"

	// CodeKey for meta.
	CodeKey = "code"

	// DurationKey for meta.
	DurationKey = "duration"
)

// WithRequestID for transport.
func WithRequestID(ctx context.Context, id meta.Value) context.Context {
	return meta.WithAttribute(ctx, RequestIDKey, id)
}

// RequestID for transport.
func RequestID(ctx context.Context) meta.Value {
	return meta.Attribute(ctx, RequestIDKey)
}
