package meta

import (
	"context"
	"unique"

	"github.com/alexfalkowski/go-service/limiter"
	"github.com/alexfalkowski/go-service/meta"
)

var (
	// UserAgentKey for meta.
	UserAgentKey = unique.Make("userAgent")

	// IPAddrKey for meta.
	IPAddrKey = unique.Make("ipAddr")

	// IPAddrKindKey for meta.
	IPAddrKindKey = unique.Make("ipAddrKind")

	// Authorization for meta.
	AuthorizationKey = unique.Make("authorization")

	// GeoLocation for meta.
	GeolocationKey = unique.Make("geoLocation")
)

// RegisterKeys for limiter.
func RegisterKeys() {
	limiter.RegisterKey("user-agent", UserAgent)
	limiter.RegisterKey("ip", IPAddr)
	limiter.RegisterKey("token", Authorization)
}

// WithUserAgent for transport.
func WithUserAgent(ctx context.Context, userAgent meta.Valuer) context.Context {
	return meta.WithAttribute(ctx, UserAgentKey.Value(), userAgent)
}

// UserAgent for transport.
func UserAgent(ctx context.Context) meta.Valuer {
	return meta.Attribute(ctx, UserAgentKey.Value())
}

// WithIPAddr for transport.
func WithIPAddr(ctx context.Context, addr meta.Valuer) context.Context {
	return meta.WithAttribute(ctx, IPAddrKey.Value(), addr)
}

// IPAddr for transport.
func IPAddr(ctx context.Context) meta.Valuer {
	return meta.Attribute(ctx, IPAddrKey.Value())
}

// WithGeolocation for transport.
func WithGeolocation(ctx context.Context, location meta.Valuer) context.Context {
	return meta.WithAttribute(ctx, GeolocationKey.Value(), location)
}

// Geolocation for transport.
func Geolocation(ctx context.Context) meta.Valuer {
	return meta.Attribute(ctx, GeolocationKey.Value())
}

// WithIPAddrKind for transport.
func WithIPAddrKind(ctx context.Context, kind meta.Valuer) context.Context {
	return meta.WithAttribute(ctx, IPAddrKindKey.Value(), kind)
}

// WithAuthorization for transport.
func WithAuthorization(ctx context.Context, auth meta.Valuer) context.Context {
	return meta.WithAttribute(ctx, AuthorizationKey.Value(), auth)
}

// Authorization for transport.
func Authorization(ctx context.Context) meta.Valuer {
	return meta.Attribute(ctx, AuthorizationKey.Value())
}
