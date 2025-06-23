package meta

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/net"
	"github.com/alexfalkowski/go-service/v2/transport/header"
	"github.com/alexfalkowski/go-service/v2/transport/strings"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// Map is an alias for metadata.MD.
type Map = metadata.MD

var (
	// Authorization is an alias for meta.Authorization.
	Authorization = meta.Authorization

	// Ignored is an alias for meta.Ignored.
	Ignored = meta.Ignored

	// NewOutgoingContext is an alias for metadata.NewOutgoingContext.
	NewOutgoingContext = metadata.NewOutgoingContext

	// Pairs is an alias for metadata.Pairs.
	Pairs = metadata.Pairs

	// WithUserID is an alias for meta.WithUserID.
	WithUserID = meta.WithUserID

	// WithAuthorization is an alias for meta.WithAuthorization.
	WithAuthorization = meta.WithAuthorization
)

// UnaryServerInterceptor for meta.
func UnaryServerInterceptor(userAgent env.UserAgent, version env.Version, generator id.Generator) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		p := info.FullMethod[1:]
		if strings.IsObservable(p) {
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
		ctx = meta.WithAuthorization(ctx, extractAuthorization(ctx, md))

		_ = grpc.SetHeader(ctx, metadata.Pairs("service-version", version.String(), "request-id", id.Value()))

		return handler(ctx, req)
	}
}

// StreamServerInterceptor for meta.
func StreamServerInterceptor(userAgent env.UserAgent, version env.Version, generator id.Generator) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		p := info.FullMethod[1:]
		if strings.IsObservable(p) {
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
		ctx = meta.WithAuthorization(ctx, extractAuthorization(ctx, md))

		wrappedStream := middleware.WrapServerStream(stream)
		wrappedStream.WrappedContext = ctx

		return handler(srv, wrappedStream)
	}
}

// UnaryClientInterceptor for meta.
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

// StreamClientInterceptor for meta.
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

func extractAuthorization(ctx context.Context, md metadata.MD) meta.Value {
	a := authorization(md)
	if strings.IsEmpty(a) {
		return meta.Blank()
	}

	_, value, err := header.ParseAuthorization(a)
	if err != nil {
		meta.WithAttribute(ctx, "authError", meta.Error(err))

		return meta.Blank()
	}

	return meta.Ignored(value)
}

func authorization(md metadata.MD) string {
	if a := md.Get("authorization"); len(a) > 0 {
		return a[0]
	}

	return ""
}

func extractGeolocation(md metadata.MD) meta.Value {
	if id := md.Get("geolocation"); len(id) > 0 {
		return meta.String(id[0])
	}

	return meta.Blank()
}
