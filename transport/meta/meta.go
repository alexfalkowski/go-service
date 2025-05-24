package meta

import (
	"context"

	"github.com/alexfalkowski/go-service/v2/meta"
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

// Value is an alias for meta.Value.
type Value = meta.Value

var (
	// String is an alias for meta.String.
	String = meta.String

	// WithAttribute is an alias for meta.WithAttribute.
	WithAttribute = meta.WithAttribute

	// Blank is an alias for meta.Blank.
	Blank = meta.Blank

	// Error is an alias for meta.Error.
	Error = meta.Error

	// Ignored is an alias for meta.Ignored.
	Ignored = meta.Ignored
)

// WithRequestID for transport.
func WithRequestID(ctx context.Context, id Value) context.Context {
	return WithAttribute(ctx, RequestIDKey, id)
}

// RequestID for transport.
func RequestID(ctx context.Context) Value {
	return meta.Attribute(ctx, RequestIDKey)
}
