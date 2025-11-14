package meta

import "github.com/alexfalkowski/go-service/v2/context"

const (
	// RequestIDKey for meta.
	RequestIDKey = "requestId"

	// SystemKey for meta.
	SystemKey = "system"

	// ServiceKey for meta.
	ServiceKey = "service"

	// MethodKey for meta.
	MethodKey = "method"

	// CodeKey for meta.
	CodeKey = "code"

	// DurationKey for meta.
	DurationKey = "duration"

	// UserAgentKey for meta.
	UserAgentKey = "userAgent"

	// UserIDKey for meta.
	UserIDKey = "userId"

	// IPAddrKey for meta.
	IPAddrKey = "ipAddr"

	// AuthorizationKey for meta.
	AuthorizationKey = "authorization"

	// GeolocationKey for meta.
	GeolocationKey = "geoLocation"
)

// WithUserAgent for meta.
func WithUserAgent(ctx context.Context, userAgent Value) context.Context {
	return WithAttribute(ctx, UserAgentKey, userAgent)
}

// UserAgent for meta.
func UserAgent(ctx context.Context) Value {
	return Attribute(ctx, UserAgentKey)
}

// WithUserID for meta.
func WithUserID(ctx context.Context, id Value) context.Context {
	return WithAttribute(ctx, UserIDKey, id)
}

// UserID for meta.
func UserID(ctx context.Context) Value {
	return Attribute(ctx, UserIDKey)
}

// WithIPAddr for meta.
func WithIPAddr(ctx context.Context, addr Value) context.Context {
	return WithAttribute(ctx, IPAddrKey, addr)
}

// IPAddr for meta.
func IPAddr(ctx context.Context) Value {
	return Attribute(ctx, IPAddrKey)
}

// WithGeolocation for meta.
func WithGeolocation(ctx context.Context, location Value) context.Context {
	return WithAttribute(ctx, GeolocationKey, location)
}

// Geolocation for meta.
func Geolocation(ctx context.Context) Value {
	return Attribute(ctx, GeolocationKey)
}

// WithIPAddrKind for meta.
func WithIPAddrKind(ctx context.Context, kind Value) context.Context {
	return WithAttribute(ctx, "ipAddrKind", kind)
}

// WithAuthorization for meta.
func WithAuthorization(ctx context.Context, auth Value) context.Context {
	return WithAttribute(ctx, AuthorizationKey, auth)
}

// Authorization for meta.
func Authorization(ctx context.Context) Value {
	return Attribute(ctx, AuthorizationKey)
}

// WithRequestID for meta.
func WithRequestID(ctx context.Context, id Value) context.Context {
	return WithAttribute(ctx, RequestIDKey, id)
}

// RequestID for meta.
func RequestID(ctx context.Context) Value {
	return Attribute(ctx, RequestIDKey)
}
