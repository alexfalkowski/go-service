package meta

import (
	"context"

	"github.com/alexfalkowski/go-service/meta"
)

const (
	// RequestIDKey for meta.
	RequestIDKey = "requestId"

	// UserAgentKey for meta.
	UserAgentKey = "userAgent"

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

	// StartTimeKey for meta.
	StartTimeKey = "startTime"

	// DeadlineKey for meta.
	DeadlineKey = "deadline"
)

// WithRequestID for transport.
func WithRequestID(ctx context.Context, id meta.Valuer) context.Context {
	return meta.WithAttribute(ctx, RequestIDKey, id)
}

// RequestID for transport.
func RequestID(ctx context.Context) meta.Valuer {
	return meta.Attribute(ctx, RequestIDKey)
}

// WithUserAgent for transport.
func WithUserAgent(ctx context.Context, userAgent meta.Valuer) context.Context {
	return meta.WithAttribute(ctx, UserAgentKey, userAgent)
}

// UserAgent for transport.
func UserAgent(ctx context.Context) meta.Valuer {
	return meta.Attribute(ctx, UserAgentKey)
}

// WithTraceID for transport.
func WithTraceID(ctx context.Context, id meta.Valuer) context.Context {
	return meta.WithAttribute(ctx, "traceId", id)
}
