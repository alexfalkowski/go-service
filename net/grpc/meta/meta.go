package meta

import (
	"fmt"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/meta"
	"google.golang.org/grpc/metadata"
)

// Map aliases gRPC metadata.MD.
//
// It represents the incoming or outgoing metadata map attached to a gRPC
// context. The alias exists so callers can work with gRPC metadata through the
// go-service import path rather than importing `google.golang.org/grpc/metadata`
// directly.
type Map = metadata.MD

// Value aliases the root go-service metadata value type.
//
// Use it when working with request-scoped attributes such as user agents,
// request IDs, authorization tokens, and IP-derived metadata through this
// package's helper functions.
type Value = meta.Value

// Pair aliases the root go-service metadata pair type.
type Pair = meta.Pair

// IPAddrKindKey is the attribute key that describes how an IP address was derived.
//
// The value typically distinguishes between peer-derived addresses and trusted
// forwarding headers.
const IPAddrKindKey = meta.IPAddrKindKey

// Attribute returns the stored request attribute for key.
//
// If no attribute is present, it returns the zero-value [Value].
func Attribute(ctx context.Context, key string) meta.Value {
	return meta.Attribute(ctx, key)
}

// NewPair creates one metadata key/value pair for batched storage updates.
func NewPair(key string, value meta.Value) Pair {
	return meta.NewPair(key, value)
}

// WithAttributes stores all provided metadata pairs on ctx.
func WithAttributes(ctx context.Context, pairs ...Pair) context.Context {
	return meta.WithAttributes(ctx, pairs...)
}

// WithRequestID creates a request ID pair for WithAttributes.
func WithRequestID(value meta.Value) Pair {
	return meta.WithRequestID(value)
}

// WithUserAgent creates a user agent pair for WithAttributes.
func WithUserAgent(value meta.Value) Pair {
	return meta.WithUserAgent(value)
}

// WithUserID creates a user ID pair for WithAttributes.
func WithUserID(value meta.Value) Pair {
	return meta.WithUserID(value)
}

// IPAddr returns the stored IP address attribute from ctx.
//
// If no value is present, it returns the zero-value meta.Value.
func IPAddr(ctx context.Context) meta.Value {
	return meta.IPAddr(ctx)
}

// WithAuthorization creates an authorization pair for WithAttributes.
//
// Authorization values often contain secrets; prefer Ignored or Redacted values if they might be exported.
func WithAuthorization(value meta.Value) Pair {
	return meta.WithAuthorization(value)
}

// Authorization returns the authorization attribute stored on ctx.
//
// If no value is present, it returns the zero-value meta.Value.
func Authorization(ctx context.Context) meta.Value {
	return meta.Authorization(ctx)
}

// AppendToOutgoingContext adds kv to the outgoing metadata in ctx.
//
// It is an alias for metadata.AppendToOutgoingContext. The kv slice must
// contain an even number of elements arranged as alternating key/value pairs.
func AppendToOutgoingContext(ctx context.Context, kv ...string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, kv...)
}

// Ignored constructs a [Value] that is retained in-context but omitted when
// rendered/exported.
//
// It is useful for sensitive metadata such as authorization-derived values.
func Ignored(value string) meta.Value {
	return meta.Ignored(value)
}

// FromOutgoingContext returns the outgoing metadata in ctx, if any.
//
// It is an alias for metadata.FromOutgoingContext.
func FromOutgoingContext(ctx context.Context) (Map, bool) {
	return metadata.FromOutgoingContext(ctx)
}

// FromIncomingContext returns the incoming metadata in ctx, if any.
//
// It is an alias for metadata.FromIncomingContext.
func FromIncomingContext(ctx context.Context) (Map, bool) {
	return metadata.FromIncomingContext(ctx)
}

// New constructs metadata from a string map.
//
// Keys are normalized using the same rules as upstream gRPC metadata handling.
func New(md map[string]string) Map {
	return metadata.New(md)
}

// NewIncomingContext attaches md as incoming metadata to ctx.
//
// It is an alias for metadata.NewIncomingContext.
func NewIncomingContext(ctx context.Context, md Map) context.Context {
	return metadata.NewIncomingContext(ctx, md)
}

// NewOutgoingContext attaches md as outgoing metadata to ctx.
func NewOutgoingContext(ctx context.Context, md Map) context.Context {
	return metadata.NewOutgoingContext(ctx, md)
}

// Pairs constructs metadata from alternating key/value arguments.
//
// The kv slice must contain an even number of elements.
func Pairs(kv ...string) Map {
	return metadata.Pairs(kv...)
}

// Redacted constructs a [Value] that renders as a mask while preserving the
// underlying value in-context.
func Redacted(value string) meta.Value {
	return meta.Redacted(value)
}

// String constructs a normal [Value] that renders as-is.
func String(value string) meta.Value {
	return meta.String(value)
}

// ToIgnored converts st to an ignored [Value] using st.String().
//
// If st is nil, it returns a blank value.
func ToIgnored(st fmt.Stringer) meta.Value {
	return meta.ToIgnored(st)
}

// ToRedacted converts st to a redacted [Value] using st.String().
//
// If st is nil, it returns a blank value.
func ToRedacted(st fmt.Stringer) meta.Value {
	return meta.ToRedacted(st)
}

// ToString converts st to a normal [Value] using st.String().
//
// If st is nil, it returns a blank value.
func ToString(st fmt.Stringer) meta.Value {
	return meta.ToString(st)
}
