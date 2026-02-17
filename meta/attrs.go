package meta

import "github.com/alexfalkowski/go-service/v2/context"

const (
	// RequestIDKey is the attribute key used for request IDs.
	RequestIDKey = "requestId"

	// SystemKey is the attribute key used for the system name.
	SystemKey = "system"

	// ServiceKey is the attribute key used for the service name.
	ServiceKey = "service"

	// MethodKey is the attribute key used for method names.
	MethodKey = "method"

	// CodeKey is the attribute key used for status codes.
	CodeKey = "code"

	// DurationKey is the attribute key used for durations.
	DurationKey = "duration"

	// UserAgentKey is the attribute key used for user agents.
	UserAgentKey = "userAgent"

	// UserIDKey is the attribute key used for user IDs.
	UserIDKey = "userId"

	// IPAddrKey is the attribute key used for IP addresses.
	IPAddrKey = "ipAddr"

	// AuthorizationKey is the attribute key used for authorization values.
	AuthorizationKey = "authorization"

	// GeolocationKey is the attribute key used for geolocation values.
	GeolocationKey = "geoLocation"
)

// WithUserAgent stores a user agent attribute on ctx.
func WithUserAgent(ctx context.Context, userAgent Value) context.Context {
	return WithAttribute(ctx, UserAgentKey, userAgent)
}

// UserAgent returns the stored user agent attribute.
func UserAgent(ctx context.Context) Value {
	return Attribute(ctx, UserAgentKey)
}

// WithUserID stores a user ID attribute on ctx.
func WithUserID(ctx context.Context, id Value) context.Context {
	return WithAttribute(ctx, UserIDKey, id)
}

// UserID returns the stored user ID attribute.
func UserID(ctx context.Context) Value {
	return Attribute(ctx, UserIDKey)
}

// WithIPAddr stores an IP address attribute on ctx.
func WithIPAddr(ctx context.Context, addr Value) context.Context {
	return WithAttribute(ctx, IPAddrKey, addr)
}

// IPAddr returns the stored IP address attribute.
func IPAddr(ctx context.Context) Value {
	return Attribute(ctx, IPAddrKey)
}

// WithGeolocation stores a geolocation attribute on ctx.
func WithGeolocation(ctx context.Context, location Value) context.Context {
	return WithAttribute(ctx, GeolocationKey, location)
}

// Geolocation returns the stored geolocation attribute.
func Geolocation(ctx context.Context) Value {
	return Attribute(ctx, GeolocationKey)
}

// WithIPAddrKind stores the IP address kind attribute on ctx.
func WithIPAddrKind(ctx context.Context, kind Value) context.Context {
	return WithAttribute(ctx, "ipAddrKind", kind)
}

// WithAuthorization stores an authorization attribute on ctx.
func WithAuthorization(ctx context.Context, auth Value) context.Context {
	return WithAttribute(ctx, AuthorizationKey, auth)
}

// Authorization returns the stored authorization attribute.
func Authorization(ctx context.Context) Value {
	return Attribute(ctx, AuthorizationKey)
}

// WithRequestID stores a request ID attribute on ctx.
func WithRequestID(ctx context.Context, id Value) context.Context {
	return WithAttribute(ctx, RequestIDKey, id)
}

// RequestID returns the stored request ID attribute.
func RequestID(ctx context.Context) Value {
	return Attribute(ctx, RequestIDKey)
}
