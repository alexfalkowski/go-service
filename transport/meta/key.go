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
)

// RegisterKeys for limiter.
func RegisterKeys() {
	limiter.RegisterKey("user-agent", UserAgent)
	limiter.RegisterKey("ip", IPAddr)
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
