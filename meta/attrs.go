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

// WithRequestID creates a request ID pair for WithAttributes.
//
// The pair stores value under RequestIDKey.
func WithRequestID(value Value) Pair {
	return NewPair(RequestIDKey, value)
}

// RequestID returns the stored request ID attribute from ctx.
//
// If no value is present, this returns the zero-value Value.
func RequestID(ctx context.Context) Value {
	return Attribute(ctx, RequestIDKey)
}

// WithSystem creates a system pair for WithAttributes.
//
// The pair stores value under SystemKey.
func WithSystem(value Value) Pair {
	return NewPair(SystemKey, value)
}

// WithService creates a service pair for WithAttributes.
//
// The pair stores value under ServiceKey.
func WithService(value Value) Pair {
	return NewPair(ServiceKey, value)
}

// WithMethod creates a method pair for WithAttributes.
//
// The pair stores value under MethodKey.
func WithMethod(value Value) Pair {
	return NewPair(MethodKey, value)
}

// WithCode creates a status code pair for WithAttributes.
//
// The pair stores value under CodeKey.
func WithCode(value Value) Pair {
	return NewPair(CodeKey, value)
}

// WithDuration creates a duration pair for WithAttributes.
//
// The pair stores value under DurationKey.
func WithDuration(value Value) Pair {
	return NewPair(DurationKey, value)
}

// WithUserAgent creates a user agent pair for WithAttributes.
//
// The pair stores value under UserAgentKey.
func WithUserAgent(value Value) Pair {
	return NewPair(UserAgentKey, value)
}

// UserAgent returns the stored user agent attribute from ctx.
//
// If no value is present, this returns the zero-value Value.
func UserAgent(ctx context.Context) Value {
	return Attribute(ctx, UserAgentKey)
}

// WithUserID creates a user ID pair for WithAttributes.
//
// The pair stores value under UserIDKey.
func WithUserID(value Value) Pair {
	return NewPair(UserIDKey, value)
}

// UserID returns the stored user ID attribute from ctx.
//
// If no value is present, this returns the zero-value Value.
func UserID(ctx context.Context) Value {
	return Attribute(ctx, UserIDKey)
}

// WithIPAddr creates an IP address pair for WithAttributes.
//
// The pair stores value under IPAddrKey.
func WithIPAddr(value Value) Pair {
	return NewPair(IPAddrKey, value)
}

// IPAddr returns the stored IP address attribute from ctx.
//
// If no value is present, this returns the zero-value Value.
func IPAddr(ctx context.Context) Value {
	return Attribute(ctx, IPAddrKey)
}

// WithIPAddrKind creates an IP address kind pair for WithAttributes.
//
// The pair stores value under IPAddrKindKey.
func WithIPAddrKind(value Value) Pair {
	return NewPair(IPAddrKindKey, value)
}

// WithAuthorization creates an authorization pair for WithAttributes.
//
// The pair stores value under AuthorizationKey. Authorization values often
// contain secrets; prefer Ignored or Redacted values if they might be exported.
func WithAuthorization(value Value) Pair {
	return NewPair(AuthorizationKey, value)
}

// Authorization returns the stored authorization attribute from ctx.
//
// If no value is present, this returns the zero-value Value.
func Authorization(ctx context.Context) Value {
	return Attribute(ctx, AuthorizationKey)
}

// WithGeolocation creates a geolocation pair for WithAttributes.
//
// The pair stores value under GeolocationKey.
func WithGeolocation(value Value) Pair {
	return NewPair(GeolocationKey, value)
}

// Geolocation returns the stored geolocation attribute from ctx.
//
// If no value is present, this returns the zero-value Value.
func Geolocation(ctx context.Context) Value {
	return Attribute(ctx, GeolocationKey)
}
