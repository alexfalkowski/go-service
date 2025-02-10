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

	// GeoLocation for meta.
	GeolocationKey = "geoLocation"
)

// RegisterKeys for limiter.
func RegisterKeys() {
	limiter.RegisterKey("user-agent", UserAgent)
	limiter.RegisterKey("ip", IPAddr)
	limiter.RegisterKey("token", Authorization)
}

// WithUserAgent for transport.
func WithUserAgent(ctx context.Context, userAgent *meta.Value) context.Context {
	return meta.WithAttribute(ctx, UserAgentKey, userAgent)
}

// UserAgent for transport.
func UserAgent(ctx context.Context) *meta.Value {
	return meta.Attribute(ctx, UserAgentKey)
}

// WithIPAddr for transport.
func WithIPAddr(ctx context.Context, addr *meta.Value) context.Context {
	return meta.WithAttribute(ctx, IPAddrKey, addr)
}

// IPAddr for transport.
func IPAddr(ctx context.Context) *meta.Value {
	return meta.Attribute(ctx, IPAddrKey)
}

// WithGeolocation for transport.
func WithGeolocation(ctx context.Context, location *meta.Value) context.Context {
	return meta.WithAttribute(ctx, GeolocationKey, location)
}

// Geolocation for transport.
func Geolocation(ctx context.Context) *meta.Value {
	return meta.Attribute(ctx, GeolocationKey)
}

// WithIPAddrKind for transport.
func WithIPAddrKind(ctx context.Context, kind *meta.Value) context.Context {
	return meta.WithAttribute(ctx, "ipAddrKind", kind)
}

// WithAuthorization for transport.
func WithAuthorization(ctx context.Context, auth *meta.Value) context.Context {
	return meta.WithAttribute(ctx, AuthorizationKey, auth)
}

// Authorization for transport.
func Authorization(ctx context.Context) *meta.Value {
	return meta.Attribute(ctx, AuthorizationKey)
}
