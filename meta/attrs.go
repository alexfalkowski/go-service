package meta

import "github.com/alexfalkowski/go-service/v2/context"

const (
	// RequestIDKey is the attribute key used for request IDs.
	//
	// This key is commonly used to correlate logs and traces for a single request across services.
	RequestIDKey = "requestId"

	// SystemKey is the attribute key used for the system name.
	//
	// This is typically used to identify the upstream system or subsystem producing/handling a request.
	SystemKey = "system"

	// ServiceKey is the attribute key used for the service name.
	//
	// This is commonly set from env.Name and used for attribution.
	ServiceKey = "service"

	// MethodKey is the attribute key used for method names.
	//
	// For example: an HTTP method name, an RPC method name, or a logical operation name.
	MethodKey = "method"

	// CodeKey is the attribute key used for status codes.
	//
	// For example: an HTTP status code or a gRPC status code.
	CodeKey = "code"

	// DurationKey is the attribute key used for durations.
	//
	// The associated value is typically rendered as a human-readable duration string.
	DurationKey = "duration"

	// UserAgentKey is the attribute key used for user agents.
	//
	// This value commonly originates from the HTTP User-Agent header.
	UserAgentKey = "userAgent"

	// UserIDKey is the attribute key used for user IDs.
	//
	// This may represent an end user, an API key identity, or a service identity depending on context.
	UserIDKey = "userId"

	// IPAddrKey is the attribute key used for IP addresses.
	//
	// This value is commonly derived from connection metadata or trusted forwarding headers.
	IPAddrKey = "ipAddr"

	// IPAddrKindKey is the attribute key used to describe how IPAddrKey was derived.
	//
	// This may be used to distinguish between direct peer IPs and values derived from proxy headers.
	IPAddrKindKey = "ipAddrKind"

	// AuthorizationKey is the attribute key used for authorization values.
	//
	// Security note: authorization values often contain secrets. Prefer storing them as Redacted or Ignored
	// values if there is any chance they will be exported to logs or headers.
	AuthorizationKey = "authorization"

	// GeolocationKey is the attribute key used for geolocation values.
	GeolocationKey = "geoLocation"
)

// WithUserAgent stores a user agent attribute on ctx under UserAgentKey.
//
// This is a convenience wrapper over WithAttribute.
func WithUserAgent(ctx context.Context, userAgent Value) context.Context {
	return WithAttribute(ctx, UserAgentKey, userAgent)
}

// UserAgent returns the stored user agent attribute from ctx.
//
// If no value is present, this returns the zero-value Value.
func UserAgent(ctx context.Context) Value {
	return Attribute(ctx, UserAgentKey)
}

// WithUserID stores a user ID attribute on ctx under UserIDKey.
//
// This is a convenience wrapper over WithAttribute.
func WithUserID(ctx context.Context, id Value) context.Context {
	return WithAttribute(ctx, UserIDKey, id)
}

// UserID returns the stored user ID attribute from ctx.
//
// If no value is present, this returns the zero-value Value.
func UserID(ctx context.Context) Value {
	return Attribute(ctx, UserIDKey)
}

// WithIPAddr stores an IP address attribute on ctx under IPAddrKey.
//
// This is a convenience wrapper over WithAttribute.
func WithIPAddr(ctx context.Context, addr Value) context.Context {
	return WithAttribute(ctx, IPAddrKey, addr)
}

// IPAddr returns the stored IP address attribute from ctx.
//
// If no value is present, this returns the zero-value Value.
func IPAddr(ctx context.Context) Value {
	return Attribute(ctx, IPAddrKey)
}

// WithGeolocation stores a geolocation attribute on ctx under GeolocationKey.
//
// This is a convenience wrapper over WithAttribute.
func WithGeolocation(ctx context.Context, location Value) context.Context {
	return WithAttribute(ctx, GeolocationKey, location)
}

// Geolocation returns the stored geolocation attribute from ctx.
//
// If no value is present, this returns the zero-value Value.
func Geolocation(ctx context.Context) Value {
	return Attribute(ctx, GeolocationKey)
}

// WithIPAddrKind stores the IP address kind attribute on ctx under IPAddrKindKey.
//
// This is a convenience wrapper over WithAttribute. It is useful for tracking whether an IP address was
// derived directly (e.g. peer IP) or indirectly (e.g. from trusted proxy headers).
func WithIPAddrKind(ctx context.Context, kind Value) context.Context {
	return WithAttribute(ctx, IPAddrKindKey, kind)
}

// WithAuthorization stores an authorization attribute on ctx under AuthorizationKey.
//
// This is a convenience wrapper over WithAttribute. Because authorization values often contain secrets,
// callers should prefer using Redacted or Ignored values unless the raw value is strictly required in-process.
func WithAuthorization(ctx context.Context, auth Value) context.Context {
	return WithAttribute(ctx, AuthorizationKey, auth)
}

// Authorization returns the stored authorization attribute from ctx.
//
// If no value is present, this returns the zero-value Value.
func Authorization(ctx context.Context) Value {
	return Attribute(ctx, AuthorizationKey)
}

// WithRequestID stores a request ID attribute on ctx under RequestIDKey.
//
// This is a convenience wrapper over WithAttribute.
func WithRequestID(ctx context.Context, id Value) context.Context {
	return WithAttribute(ctx, RequestIDKey, id)
}

// RequestID returns the stored request ID attribute from ctx.
//
// If no value is present, this returns the zero-value Value.
func RequestID(ctx context.Context) Value {
	return Attribute(ctx, RequestIDKey)
}
