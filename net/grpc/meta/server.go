package meta

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/net"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/net/grpc/method"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/alexfalkowski/go-service/v2/net/header"
	"github.com/alexfalkowski/go-service/v2/slices"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/telemetry/attributes"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// UnaryServerInterceptor returns a gRPC unary server interceptor that extracts metadata into the context.
//
// Requests with ignorable unary methods bypass extraction.
//
// For non-ignored methods, the interceptor:
//
//   - copies incoming metadata from the request context
//   - stores "grpc" as the ignored transport attribute
//   - resolves "user-agent" and "request-id", preferring existing context
//     attributes and then incoming metadata values
//   - stores the RPC full method name as the ignored service-method attribute
//   - derives IP address information from trusted forwarding headers or, if
//     absent, from the gRPC peer address
//   - parses the "authorization" header into the request attribute model
//   - stores "geolocation" as ignored when present
//   - sets response header metadata including "service-version" and
//     "request-id"
//
// If the Authorization header is present but invalid, the interceptor returns a
// [codes.InvalidArgument] gRPC status error.
func UnaryServerInterceptor(policy *method.Policy, userAgent env.UserAgent, version env.Version, generator id.Generator) grpc.UnaryServerInterceptor {
	serviceVersion := version.String()

	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if policy.IsOperation(info.FullMethod) {
			return handler(ctx, req)
		}

		ua := serverUserAgent(ctx, userAgent)
		id := serverRequestID(ctx, generator)

		ipKind, ipAddr := serverIPAddr(ctx)
		geolocation := serverGeolocation(ctx)

		auth, err := serverAuthorization(ctx)
		if err != nil {
			return nil, status.SafeError(codes.InvalidArgument, err)
		}
		ctx = meta.WithAttributes(ctx,
			meta.WithTransport(meta.Ignored("grpc")),
			meta.WithUserAgent(ua),
			meta.WithRequestID(id),
			meta.WithServiceMethod(meta.Ignored(info.FullMethod)),
			meta.WithIPAddr(ipAddr),
			meta.WithIPAddrKind(ipKind),
			meta.WithGeolocation(geolocation),
			meta.WithAuthorization(auth),
		)
		stampSpan(ctx)

		_ = grpc.SetHeader(ctx, serverResponseHeaders(serviceVersion, id.Value()))

		return handler(ctx, req)
	}
}

// StreamServerInterceptor returns a gRPC stream server interceptor that extracts metadata into the stream context.
//
// Stream methods always run metadata extraction. This keeps long-lived operation streams, such as health Watch,
// visible to downstream stream interceptors that need request metadata for resource controls.
//
// The interceptor performs the same metadata-to-context projection as [UnaryServerInterceptor], but applies it to
// the wrapped stream context, stores "grpc" as the ignored transport attribute, and emits response headers
// through the stream API.
func StreamServerInterceptor(policy *method.Policy, userAgent env.UserAgent, version env.Version, generator id.Generator) grpc.StreamServerInterceptor {
	serviceVersion := version.String()

	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := stream.Context()
		ua := serverUserAgent(ctx, userAgent)

		id := serverRequestID(ctx, generator)

		ipKind, ipAddr := serverIPAddr(ctx)
		geolocation := serverGeolocation(ctx)

		// Operation streams still need metadata for limiting, but they keep the same auth bypass as unary
		// operation methods.
		var auth meta.Value
		if policy.IsOperation(info.FullMethod) {
			auth = meta.Blank()
		} else {
			a, err := serverAuthorization(ctx)
			if err != nil {
				return status.SafeError(codes.InvalidArgument, err)
			}
			auth = a
		}
		ctx = meta.WithAttributes(ctx,
			meta.WithTransport(meta.Ignored("grpc")),
			meta.WithUserAgent(ua),
			meta.WithRequestID(id),
			meta.WithServiceMethod(meta.Ignored(info.FullMethod)),
			meta.WithIPAddr(ipAddr),
			meta.WithIPAddrKind(ipKind),
			meta.WithGeolocation(geolocation),
			meta.WithAuthorization(auth),
		)
		stampSpan(ctx)

		_ = stream.SetHeader(serverResponseHeaders(serviceVersion, id.Value()))

		wrappedStream := middleware.WrapServerStream(stream)
		wrappedStream.WrappedContext = ctx

		return handler(srv, wrappedStream)
	}
}

// stampSpan copies the request metadata in ctx onto the active span so the
// server/root span carries the same request/service context as the logs.
//
// The tracer's metadata span processor cannot cover the server span, because
// its context has no metadata yet when the span starts; child spans (database,
// cache, ...) are stamped by that processor instead.
func stampSpan(ctx context.Context) {
	attributes.Record(ctx, attributes.Strings(meta.CamelStrings(ctx, meta.NoPrefix))...)
}

func serverResponseHeaders(serviceVersion, requestID string) Map {
	// Clip caps each metadata value at one element so later appends allocate
	// instead of overwriting the neighboring value in this backing array.
	values := [...]string{serviceVersion, requestID}
	return Map{
		"service-version": slices.Clip(values[0:1]),
		"request-id":      slices.Clip(values[1:2]),
	}
}

func serverIPAddr(ctx context.Context) (meta.Value, meta.Value) {
	for _, h := range header.ForwardedIPs {
		k := h.GRPC
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

	return peerKind, meta.String(serverPeerIPAddr(peer.Addr))
}

func serverPeerIPAddr(addr net.Addr) string {
	switch addr := addr.(type) {
	case *net.TCPAddr:
		return addr.IP.String()
	case *net.UDPAddr:
		return addr.IP.String()
	default:
		return net.Host(addr.String())
	}
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

func serverRequestID(ctx context.Context, generator id.Generator) meta.Value {
	if id := meta.RequestID(ctx); !id.IsEmpty() {
		return id
	}
	if id := serverValue(ctx, "request-id"); !strings.IsEmpty(id) {
		return meta.String(id)
	}

	return meta.String(generator.Generate())
}

func serverAuthorization(ctx context.Context) (meta.Value, error) {
	a := serverValue(ctx, "authorization")
	if strings.IsEmpty(a) {
		return meta.Blank(), nil
	}

	value, err := header.ParseBearer(a)
	if err != nil {
		return meta.Blank(), err
	}

	return meta.Ignored(value), nil
}

func serverGeolocation(ctx context.Context) meta.Value {
	if id := serverValue(ctx, "geolocation"); !strings.IsEmpty(id) {
		return meta.Ignored(id)
	}
	return meta.Blank()
}

func serverValue(ctx context.Context, key string) string {
	if values := metadata.ValueFromIncomingContext(ctx, key); len(values) > 0 {
		return values[0]
	}
	return strings.Empty
}
