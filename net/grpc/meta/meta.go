package meta

import (
	"fmt"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/net"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/alexfalkowski/go-service/v2/net/grpc/strings"
	"github.com/alexfalkowski/go-service/v2/net/header"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
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

// UnaryServerInterceptor returns a gRPC unary server interceptor that extracts metadata into the context.
//
// Requests with ignorable methods bypass extraction.
//
// For non-ignored methods, the interceptor:
//
//   - copies incoming metadata from the request context
//   - resolves "user-agent" and "request-id", preferring existing context
//     attributes and then incoming metadata values
//   - derives IP address information from trusted forwarding headers or, if
//     absent, from the gRPC peer address
//   - parses the "authorization" header into the request attribute model
//   - stores "geolocation" when present
//   - sets response header metadata including "service-version" and
//     "request-id"
//
// If the Authorization header is present but invalid, the interceptor returns a
// `codes.InvalidArgument` gRPC status error.
func UnaryServerInterceptor(userAgent env.UserAgent, version env.Version, generator id.Generator) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if strings.IsIgnorable(info.FullMethod) {
			return handler(ctx, req)
		}

		ua := serverUserAgent(ctx, userAgent)
		id := serverRequestID(ctx, generator)

		kind, ip := serverIPAddr(ctx)
		geolocation := serverGeolocation(ctx)

		auth, err := serverAuthorization(ctx)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		ctx = meta.WithAttributes(ctx,
			meta.WithUserAgent(ua),
			meta.WithRequestID(id),
			meta.WithIPAddr(ip),
			meta.WithIPAddrKind(kind),
			meta.WithGeolocation(geolocation),
			meta.WithAuthorization(auth),
		)

		_ = grpc.SetHeader(ctx, Pairs("service-version", version.String(), "request-id", id.Value()))

		return handler(ctx, req)
	}
}

// StreamServerInterceptor returns a gRPC stream server interceptor that extracts metadata into the stream context.
//
// Requests with ignorable methods bypass extraction.
//
// For non-ignored methods, the interceptor performs the same metadata-to-context
// projection as [UnaryServerInterceptor], but applies it to the wrapped stream
// context and emits response headers through the stream API.
func StreamServerInterceptor(userAgent env.UserAgent, version env.Version, generator id.Generator) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if strings.IsIgnorable(info.FullMethod) {
			return handler(srv, stream)
		}

		_ = stream.SetHeader(Pairs("service-version", version.String()))

		ctx := stream.Context()
		ua := serverUserAgent(ctx, userAgent)

		id := serverRequestID(ctx, generator)
		_ = stream.SetHeader(Pairs("request-id", id.Value()))

		kind, ip := serverIPAddr(ctx)
		geolocation := serverGeolocation(ctx)

		auth, err := serverAuthorization(ctx)
		if err != nil {
			return status.Error(codes.InvalidArgument, err.Error())
		}
		ctx = meta.WithAttributes(ctx,
			meta.WithUserAgent(ua),
			meta.WithRequestID(id),
			meta.WithIPAddr(ip),
			meta.WithIPAddrKind(kind),
			meta.WithGeolocation(geolocation),
			meta.WithAuthorization(auth),
		)

		wrappedStream := middleware.WrapServerStream(stream)
		wrappedStream.WrappedContext = ctx

		return handler(srv, wrappedStream)
	}
}

// UnaryClientInterceptor returns a gRPC unary client interceptor that injects metadata into outgoing requests.
//
// It ensures "user-agent" and "request-id" are present in outgoing metadata,
// preferring values already present in the context or outgoing metadata, and
// stores the chosen values back into the context.
//
// Existing outgoing metadata values for these keys are replaced so repeated
// interceptor invocation does not accumulate duplicates or preserve stale
// values ahead of the resolved value.
func UnaryClientInterceptor(userAgent env.UserAgent, generator id.Generator) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, fullMethod string, req, resp any, conn *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		md := ExtractOutgoing(ctx)

		ua := clientUserAgent(ctx, md, userAgent)
		id := clientRequestID(ctx, generator, md)

		md.Set("user-agent", ua.Value())
		md.Set("request-id", id.Value())

		ctx = meta.WithAttributes(ctx,
			meta.WithUserAgent(ua),
			meta.WithRequestID(id),
		)
		ctx = NewOutgoingContext(ctx, md)
		return invoker(ctx, fullMethod, req, resp, conn, opts...)
	}
}

// StreamClientInterceptor returns a gRPC stream client interceptor that injects metadata into outgoing requests.
//
// It ensures "user-agent" and "request-id" are present in outgoing metadata,
// preferring values already present in the context or outgoing metadata, and
// stores the chosen values back into the context.
//
// Existing outgoing metadata values for these keys are replaced so repeated
// interceptor invocation does not accumulate duplicates or preserve stale
// values ahead of the resolved value.
func StreamClientInterceptor(userAgent env.UserAgent, generator id.Generator) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, conn *grpc.ClientConn, fullMethod string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		md := ExtractOutgoing(ctx)

		ua := clientUserAgent(ctx, md, userAgent)
		id := clientRequestID(ctx, generator, md)

		md.Set("user-agent", ua.Value())
		md.Set("request-id", id.Value())

		ctx = meta.WithAttributes(ctx,
			meta.WithUserAgent(ua),
			meta.WithRequestID(id),
		)
		ctx = NewOutgoingContext(ctx, md)
		return streamer(ctx, desc, conn, fullMethod, opts...)
	}
}

func serverIPAddr(ctx context.Context) (meta.Value, meta.Value) {
	headers := []string{"x-real-ip", "cf-connecting-ip", "true-client-ip", "x-forwarded-for"}
	for _, k := range headers {
		if f := serverValue(ctx, k); !strings.IsEmpty(f) {
			ip, _, _ := strings.Cut(f, ",")

			return meta.String(k), meta.String(ip)
		}
	}

	peerKind := meta.String("peer")
	peer, ok := peer.FromContext(ctx)
	if !ok || peer == nil || peer.Addr == nil {
		return peerKind, meta.Blank()
	}

	return peerKind, meta.String(net.Host(peer.Addr.String()))
}

func serverUserAgent(ctx context.Context, userAgent env.UserAgent) meta.Value {
	if ua := meta.UserAgent(ctx); !ua.IsEmpty() {
		return ua
	}
	if ua := serverValue(ctx, "user-agent"); !strings.IsEmpty(ua) {
		return meta.String(ua)
	}

	return meta.String(userAgent.String())
}

func clientUserAgent(ctx context.Context, md metadata.MD, userAgent env.UserAgent) meta.Value {
	if ua := meta.UserAgent(ctx); !ua.IsEmpty() {
		return ua
	}
	if ua := md.Get("user-agent"); len(ua) > 0 {
		return meta.String(ua[0])
	}

	return meta.String(userAgent.String())
}

func serverRequestID(ctx context.Context, generator id.Generator) meta.Value {
	if id := meta.RequestID(ctx); !id.IsEmpty() {
		return id
	}
	if id := serverValue(ctx, "request-id"); !strings.IsEmpty(id) {
		return meta.String(id)
	}

	return meta.String(generator.Generate())
}

func clientRequestID(ctx context.Context, generator id.Generator, md metadata.MD) meta.Value {
	if id := meta.RequestID(ctx); !id.IsEmpty() {
		return id
	}
	if id := md.Get("request-id"); len(id) > 0 {
		return meta.String(id[0])
	}

	return meta.String(generator.Generate())
}

func serverAuthorization(ctx context.Context) (meta.Value, error) {
	a := serverValue(ctx, "authorization")
	if strings.IsEmpty(a) {
		return meta.Blank(), nil
	}

	_, value, err := header.ParseAuthorization(a)
	if err != nil {
		return meta.Blank(), err
	}

	return meta.Ignored(value), nil
}

func serverGeolocation(ctx context.Context) meta.Value {
	if id := serverValue(ctx, "geolocation"); !strings.IsEmpty(id) {
		return meta.String(id)
	}
	return meta.Blank()
}

func serverValue(ctx context.Context, key string) string {
	if values := metadata.ValueFromIncomingContext(ctx, key); len(values) > 0 {
		return values[0]
	}
	return strings.Empty
}
