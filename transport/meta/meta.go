package meta

import (
	"context"
	"unique"

	"github.com/alexfalkowski/go-service/meta"
)

var (
	// RequestIDKey for meta.
	RequestIDKey = unique.Make("requestId")

	// ServiceKey for meta.
	ServiceKey = unique.Make("service")

	// PathKey for meta.
	PathKey = unique.Make("path")

	// MethodKey for meta.
	MethodKey = unique.Make("method")

	// CodeKey for meta.
	CodeKey = unique.Make("code")

	// DurationKey for meta.
	DurationKey = unique.Make("duration")
)

// WithRequestID for transport.
func WithRequestID(ctx context.Context, id meta.Valuer) context.Context {
	return meta.WithAttribute(ctx, RequestIDKey.Value(), id)
}

// RequestID for transport.
func RequestID(ctx context.Context) meta.Valuer {
	return meta.Attribute(ctx, RequestIDKey.Value())
}
