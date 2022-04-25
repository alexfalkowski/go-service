package meta

import (
	"context"

	"github.com/alexfalkowski/go-service/meta"
)

const (
	// RequestIDKey for meta.
	RequestIDKey = "transport.request_id"

	// UserAgentKey for meta.
	UserAgentKey = "transport.user_agent"

	// RemoteAddressKey for meta.
	RemoteAddressKey = "transport.remote_address"
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
func WithUserAgent(ctx context.Context, userAgent string) context.Context {
	return meta.WithAttribute(ctx, UserAgentKey, userAgent)
}

// UserAgent for transport.
func UserAgent(ctx context.Context) string {
	return meta.Attribute(ctx, UserAgentKey)
}

// WithRemoteAddress for transport.
func WithRemoteAddress(ctx context.Context, remoteAddress string) context.Context {
	return meta.WithAttribute(ctx, RemoteAddressKey, remoteAddress)
}

// RemoteAddress for transport.
func RemoteAddress(ctx context.Context) string {
	return meta.Attribute(ctx, RemoteAddressKey)
}
