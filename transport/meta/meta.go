package meta

import (
	"context"

	"github.com/alexfalkowski/go-service/v2/limiter"
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

	// UserAgentKey for meta.
	UserAgentKey = "userAgent"

	// UserIDKey for token.
	UserIDKey = "userId"

	// IPAddrKey for meta.
	IPAddrKey = "ipAddr"

	// Authorization for meta.
	AuthorizationKey = "authorization"

	// GeoLocation for meta.
	GeolocationKey = "geoLocation"
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

// RegisterKeys for limiter.
func RegisterKeys() {
	limiter.RegisterKey("user-agent", UserAgent)
	limiter.RegisterKey("ip", IPAddr)
	limiter.RegisterKey("token", Authorization)
}

// WithUserAgent for transport.
func WithUserAgent(ctx context.Context, userAgent meta.Value) context.Context {
	return meta.WithAttribute(ctx, UserAgentKey, userAgent)
}

// UserAgent for transport.
func UserAgent(ctx context.Context) meta.Value {
	return meta.Attribute(ctx, UserAgentKey)
}

// WithUserID for token.
func WithUserID(ctx context.Context, id meta.Value) context.Context {
	return meta.WithAttribute(ctx, UserIDKey, id)
}

// UserID for token.
func UserID(ctx context.Context) meta.Value {
	return meta.Attribute(ctx, UserIDKey)
}

// WithIPAddr for transport.
func WithIPAddr(ctx context.Context, addr meta.Value) context.Context {
	return meta.WithAttribute(ctx, IPAddrKey, addr)
}

// IPAddr for transport.
func IPAddr(ctx context.Context) meta.Value {
	return meta.Attribute(ctx, IPAddrKey)
}

// WithGeolocation for transport.
func WithGeolocation(ctx context.Context, location meta.Value) context.Context {
	return meta.WithAttribute(ctx, GeolocationKey, location)
}

// Geolocation for transport.
func Geolocation(ctx context.Context) meta.Value {
	return meta.Attribute(ctx, GeolocationKey)
}

// WithIPAddrKind for transport.
func WithIPAddrKind(ctx context.Context, kind meta.Value) context.Context {
	return meta.WithAttribute(ctx, "ipAddrKind", kind)
}

// WithAuthorization for transport.
func WithAuthorization(ctx context.Context, auth meta.Value) context.Context {
	return meta.WithAttribute(ctx, AuthorizationKey, auth)
}

// Authorization for transport.
func Authorization(ctx context.Context) meta.Value {
	return meta.Attribute(ctx, AuthorizationKey)
}

// WithRequestID for transport.
func WithRequestID(ctx context.Context, id Value) context.Context {
	return WithAttribute(ctx, RequestIDKey, id)
}

// RequestID for transport.
func RequestID(ctx context.Context) Value {
	return meta.Attribute(ctx, RequestIDKey)
}
