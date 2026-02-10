package meta

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/net"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/alexfalkowski/go-service/v2/transport/header"
	"github.com/alexfalkowski/go-service/v2/transport/strings"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// Map is an alias for metadata.MD.
type Map = metadata.MD

// Authorization is an alias for meta.Authorization.
func Authorization(ctx context.Context) meta.Value {
	return meta.Authorization(ctx)
}

// Ignored is an alias for meta.Ignored.
func Ignored(value string) meta.Value {
	return meta.Ignored(value)
}

// NewOutgoingContext is an alias for metadata.NewOutgoingContext.
func NewOutgoingContext(ctx context.Context, md Map) context.Context {
	return metadata.NewOutgoingContext(ctx, md)
}

// Pairs is an alias for metadata.Pairs.
func Pairs(kv ...string) Map {
	return metadata.Pairs(kv...)
}

// WithUserID is an alias for meta.WithUserID.
func WithUserID(ctx context.Context, id meta.Value) context.Context {
	return meta.WithUserID(ctx, id)
}

// WithAuthorization is an alias for meta.WithAuthorization.
func WithAuthorization(ctx context.Context, auth meta.Value) context.Context {
	return meta.WithAuthorization(ctx, auth)
}

// UnaryServerInterceptor returns a gRPC unary server interceptor that extracts metadata into the context.
//
// Requests with ignorable methods bypass extraction.
// It extracts metadata from incoming headers (when present) and stores it into the context, including:
// "user-agent", "request-id", "authorization", and "geolocation", along with IP address and its source kind.
// If the Authorization header is present but invalid, it returns InvalidArgument.
//
// It also sets response header metadata including "service-version" and "request-id".
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

		_ = grpc.SetHeader(ctx, metadata.Pairs("service-version", version.String(), "request-id", id.Value()))

		return handler(ctx, req)
	}
}

// StreamServerInterceptor returns a gRPC stream server interceptor that extracts metadata into the stream context.
//
// Requests with ignorable methods bypass extraction.
// It extracts metadata from incoming headers (when present) and stores it into the stream context, including:
// "user-agent", "request-id", "authorization", and "geolocation", along with IP address and its source kind.
//
// It sets response header metadata including "service-version" and "request-id".
func StreamServerInterceptor(userAgent env.UserAgent, version env.Version, generator id.Generator) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if strings.IsIgnorable(info.FullMethod) {
			return handler(srv, stream)
		}

		_ = stream.SetHeader(metadata.Pairs("service-version", version.String()))

		ctx := stream.Context()
		md := ExtractIncoming(ctx)
		ctx = meta.WithUserAgent(ctx, extractUserAgent(ctx, md, userAgent))

		id := extractRequestID(ctx, generator, md)
		_ = stream.SetHeader(metadata.Pairs("request-id", id.Value()))

		ctx = meta.WithRequestID(ctx, id)

		kind, ip := extractIPAddr(ctx, md)
		ctx = meta.WithIPAddr(ctx, ip)
		ctx = meta.WithIPAddrKind(ctx, kind)

		ctx = meta.WithGeolocation(ctx, extractGeolocation(md))

		auth, err := extractAuthorization(md)
		if err != nil {
			return err
		}
		ctx = meta.WithAuthorization(ctx, auth)

		wrappedStream := middleware.WrapServerStream(stream)
		wrappedStream.WrappedContext = ctx

		return handler(srv, wrappedStream)
	}
}

// UnaryClientInterceptor returns a gRPC unary client interceptor that injects metadata into outgoing requests.
//
// It ensures "user-agent" and "request-id" are present in outgoing metadata, preferring values already
// present in the context or outgoing metadata, and stores the chosen values back into the context.
func UnaryClientInterceptor(userAgent env.UserAgent, generator id.Generator) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, fullMethod string, req, resp any, conn *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		md := ExtractOutgoing(ctx)

		ua := extractUserAgent(ctx, md, userAgent)
		ctx = meta.WithUserAgent(ctx, ua)
		md.Append("user-agent", ua.Value())

		id := extractRequestID(ctx, generator, md)
		ctx = meta.WithRequestID(ctx, id)
		md.Append("request-id", id.Value())

		ctx = metadata.NewOutgoingContext(ctx, md)
		return invoker(ctx, fullMethod, req, resp, conn, opts...)
	}
}

// StreamClientInterceptor returns a gRPC stream client interceptor that injects metadata into outgoing requests.
//
// It ensures "user-agent" and "request-id" are present in outgoing metadata, preferring values already
// present in the context or outgoing metadata, and stores the chosen values back into the context.
func StreamClientInterceptor(userAgent env.UserAgent, generator id.Generator) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, conn *grpc.ClientConn, fullMethod string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		md := ExtractOutgoing(ctx)

		ua := extractUserAgent(ctx, md, userAgent)
		ctx = meta.WithUserAgent(ctx, ua)
		md.Append("user-agent", ua.Value())

		id := extractRequestID(ctx, generator, md)
		ctx = meta.WithRequestID(ctx, id)
		md.Append("request-id", id.Value())

		ctx = metadata.NewOutgoingContext(ctx, md)
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
	peer, _ := peer.FromContext(ctx)
	addr := peer.Addr.String()

	return peerKind, meta.String(net.Host(addr))
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
