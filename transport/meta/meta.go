package meta

import (
	"context"
	"fmt"

	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/os"
)

var (
	// RequestIDKey for meta.
	// nolint:gochecknoglobals
	RequestIDKey = fmt.Sprintf("%s.request_id", os.ExecutableName())

	// UserAgentKey for meta.
	// nolint:gochecknoglobals
	UserAgentKey = fmt.Sprintf("%s.user_agent", os.ExecutableName())
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
