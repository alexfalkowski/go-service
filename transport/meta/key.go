package meta

import (
	"context"

	"github.com/alexfalkowski/go-service/limiter"
	"github.com/alexfalkowski/go-service/meta"
)

const (
	// UserAgentKey for meta.
	UserAgentKey = "userAgent"

	// IPAddrKey for meta.
	IPAddrKey = "ipAddr"

	// Authorization for meta.
	AuthorizationKey = "authorization"
)

// RegisterKeys for limiter.
func RegisterKeys() {
	limiter.RegisterKey("user-agent", UserAgent)
	limiter.RegisterKey("ip", IPAddr)
	limiter.RegisterKey("token", Authorization)
}

// WithUserAgent for transport.
func WithUserAgent(ctx context.Context, userAgent meta.Valuer) context.Context {
	return meta.WithAttribute(ctx, UserAgentKey, userAgent)
}

// UserAgent for transport.
func UserAgent(ctx context.Context) meta.Valuer {
	return meta.Attribute(ctx, UserAgentKey)
}

// WithIPAddr for transport.
func WithIPAddr(ctx context.Context, addr meta.Valuer) context.Context {
	return meta.WithAttribute(ctx, IPAddrKey, addr)
}

// IPAddr for transport.
func IPAddr(ctx context.Context) meta.Valuer {
	return meta.Attribute(ctx, IPAddrKey)
}

// WithIPAddrKind for transport.
func WithIPAddrKind(ctx context.Context, kind meta.Valuer) context.Context {
	return meta.WithAttribute(ctx, "ipAddrKind", kind)
}

// WithAuthorization for transport.
func WithAuthorization(ctx context.Context, auth meta.Valuer) context.Context {
	return meta.WithAttribute(ctx, AuthorizationKey, auth)
}

// Authorization for transport.
func Authorization(ctx context.Context) meta.Valuer {
	return meta.Attribute(ctx, AuthorizationKey)
}
