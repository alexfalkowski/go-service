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

// IPAddrKindKey is the attribute key that describes how an IP address was derived.
//
// The value typically distinguishes between peer-derived addresses and trusted
// forwarding headers.
const IPAddrKindKey = meta.IPAddrKindKey

// Authorization returns the authorization attribute stored on ctx.
//
// It forwards to the root `meta.Authorization` helper so callers can keep gRPC
// metadata and request attribute access under one import path.
func Authorization(ctx context.Context) meta.Value {
	return meta.Authorization(ctx)
}

// Attribute returns the stored request attribute for key.
//
// If no attribute is present, it returns the zero-value [Value]. It forwards to
// the root `meta.Attribute` helper.
func Attribute(ctx context.Context, key string) meta.Value {
	return meta.Attribute(ctx, key)
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
// It forwards to the root `meta.Ignored` helper and is useful for sensitive
// metadata such as authorization-derived values.
func Ignored(value string) meta.Value {
	return meta.Ignored(value)
}

// IPAddr returns the stored IP address attribute from ctx.
//
// It forwards to the root `meta.IPAddr` helper.
func IPAddr(ctx context.Context) meta.Value {
	return meta.IPAddr(ctx)
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
// It forwards to `metadata.New`. Keys are normalized using the same rules as
// upstream gRPC metadata handling.
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
//
// It forwards to `metadata.NewOutgoingContext`.
func NewOutgoingContext(ctx context.Context, md Map) context.Context {
	return metadata.NewOutgoingContext(ctx, md)
}

// Pairs constructs metadata from alternating key/value arguments.
//
// It forwards to `metadata.Pairs`. The kv slice must contain an even number of
// elements.
func Pairs(kv ...string) Map {
	return metadata.Pairs(kv...)
}

// Redacted constructs a [Value] that renders as a mask while preserving the
// underlying value in-context.
//
// It forwards to the root `meta.Redacted` helper.
func Redacted(value string) meta.Value {
	return meta.Redacted(value)
}

// String constructs a normal [Value] that renders as-is.
//
// It forwards to the root `meta.String` helper.
func String(value string) meta.Value {
	return meta.String(value)
}

// ToIgnored converts st to an ignored [Value] using st.String().
//
// If st is nil, it returns a blank value. It forwards to the root
// `meta.ToIgnored` helper.
func ToIgnored(st fmt.Stringer) meta.Value {
	return meta.ToIgnored(st)
}

// ToRedacted converts st to a redacted [Value] using st.String().
//
// If st is nil, it returns a blank value. It forwards to the root
// `meta.ToRedacted` helper.
func ToRedacted(st fmt.Stringer) meta.Value {
	return meta.ToRedacted(st)
}

// ToString converts st to a normal [Value] using st.String().
//
// If st is nil, it returns a blank value. It forwards to the root
// `meta.ToString` helper.
func ToString(st fmt.Stringer) meta.Value {
	return meta.ToString(st)
}

// WithAttribute stores key/value on ctx as a request-scoped attribute.
//
// It forwards to the root `meta.WithAttribute` helper and is primarily useful
// when tests or interceptors need to seed context values before gRPC transport
// processing begins.
func WithAttribute(ctx context.Context, key string, value meta.Value) context.Context {
	return meta.WithAttribute(ctx, key, value)
}

// WithRequestID stores a request ID attribute on ctx.
//
// It forwards to the root `meta.WithRequestID` helper.
func WithRequestID(ctx context.Context, id meta.Value) context.Context {
	return meta.WithRequestID(ctx, id)
}

// WithUserID stores a user ID attribute on ctx.
//
// It forwards to the root `meta.WithUserID` helper.
func WithUserID(ctx context.Context, id meta.Value) context.Context {
	return meta.WithUserID(ctx, id)
}

// WithUserAgent stores a user-agent attribute on ctx.
//
// It forwards to the root `meta.WithUserAgent` helper.
func WithUserAgent(ctx context.Context, userAgent meta.Value) context.Context {
	return meta.WithUserAgent(ctx, userAgent)
}

// WithAuthorization stores an authorization attribute on ctx.
//
// It forwards to the root `meta.WithAuthorization` helper.
func WithAuthorization(ctx context.Context, auth meta.Value) context.Context {
	return meta.WithAuthorization(ctx, auth)
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

		md := ExtractIncoming(ctx)

		ctx = meta.WithUserAgent(ctx, extractUserAgent(ctx, md, userAgent))

		id := extractRequestID(ctx, generator, md)
		ctx = meta.WithRequestID(ctx, id)

		kind, ip := extractIPAddr(ctx, md)
		ctx = meta.WithIPAddr(ctx, ip)
		ctx = meta.WithIPAddrKind(ctx, kind)

		ctx = meta.WithGeolocation(ctx, extractGeolocation(md))

		auth, err := extractAuthorization(md)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		ctx = meta.WithAuthorization(ctx, auth)

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
		md := ExtractIncoming(ctx)
		ctx = meta.WithUserAgent(ctx, extractUserAgent(ctx, md, userAgent))

		id := extractRequestID(ctx, generator, md)
		_ = stream.SetHeader(Pairs("request-id", id.Value()))

		ctx = meta.WithRequestID(ctx, id)

		kind, ip := extractIPAddr(ctx, md)
		ctx = meta.WithIPAddr(ctx, ip)
		ctx = meta.WithIPAddrKind(ctx, kind)

		ctx = meta.WithGeolocation(ctx, extractGeolocation(md))

		auth, err := extractAuthorization(md)
		if err != nil {
			return status.Error(codes.InvalidArgument, err.Error())
		}
		ctx = meta.WithAuthorization(ctx, auth)

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

		ua := extractUserAgent(ctx, md, userAgent)
		ctx = meta.WithUserAgent(ctx, ua)
		md.Set("user-agent", ua.Value())

		id := extractRequestID(ctx, generator, md)
		ctx = meta.WithRequestID(ctx, id)
		md.Set("request-id", id.Value())

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

		ua := extractUserAgent(ctx, md, userAgent)
		ctx = meta.WithUserAgent(ctx, ua)
		md.Set("user-agent", ua.Value())

		id := extractRequestID(ctx, generator, md)
		ctx = meta.WithRequestID(ctx, id)
		md.Set("request-id", id.Value())

		ctx = NewOutgoingContext(ctx, md)
		return streamer(ctx, desc, conn, fullMethod, opts...)
	}
}

func extractIPAddr(ctx context.Context, md metadata.MD) (meta.Value, meta.Value) {
	headers := []string{"x-real-ip", "cf-connecting-ip", "true-client-ip", "x-forwarded-for"}
	for _, k := range headers {
		if f := md.Get(k); len(f) > 0 {
			ip, _, _ := strings.Cut(f[0], ",")

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

func extractUserAgent(ctx context.Context, md metadata.MD, userAgent env.UserAgent) meta.Value {
	if ua := meta.UserAgent(ctx); !ua.IsEmpty() {
		return ua
	}
	if ua := md.Get("user-agent"); len(ua) > 0 {
		return meta.String(ua[0])
	}

	return meta.String(userAgent.String())
}

func extractRequestID(ctx context.Context, generator id.Generator, md metadata.MD) meta.Value {
	if id := meta.RequestID(ctx); !id.IsEmpty() {
		return id
	}
	if id := md.Get("request-id"); len(id) > 0 {
		return meta.String(id[0])
	}

	return meta.String(generator.Generate())
}

func extractAuthorization(md metadata.MD) (meta.Value, error) {
	a := authorization(md)
	if strings.IsEmpty(a) {
		return meta.Blank(), nil
	}

	_, value, err := header.ParseAuthorization(a)
	if err != nil {
		return meta.Blank(), err
	}

	return meta.Ignored(value), nil
}

func authorization(md metadata.MD) string {
	if a := md.Get("authorization"); len(a) > 0 {
		return a[0]
	}
	return strings.Empty
}

func extractGeolocation(md metadata.MD) meta.Value {
	if id := md.Get("geolocation"); len(id) > 0 {
		return meta.String(id[0])
	}
	return meta.Blank()
}
