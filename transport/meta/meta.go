package meta

import (
	"context"

	"github.com/alexfalkowski/go-service/meta"
)

const (
	// RequestIDKey for meta.
	RequestIDKey = "app.request_id"

	// UserAgentKey for meta.
	UserAgentKey = "app.user_agent"
)

// WithRequestID for transport.
func WithRequestID(ctx context.Context, id string) context.Context {
	return meta.WithAttribute(ctx, RequestIDKey, id)
}

// RequestID for transport.
func RequestID(ctx context.Context) string {
	return meta.Attribute(ctx, RequestIDKey)
}

// WithUserAgent for transport.
func WithUserAgent(ctx context.Context, id string) context.Context {
	return meta.WithAttribute(ctx, UserAgentKey, id)
}

// UserAgent for transport.
func UserAgent(ctx context.Context) string {
	return meta.Attribute(ctx, UserAgentKey)
}
