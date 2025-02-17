package meta

import (
	"context"
	"path"
	"strings"

	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/id"
	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/net"
	"github.com/alexfalkowski/go-service/transport/header"
	m "github.com/alexfalkowski/go-service/transport/meta"
	ts "github.com/alexfalkowski/go-service/transport/strings"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// UnaryServerInterceptor for meta.
func UnaryServerInterceptor(userAgent env.UserAgent, version env.Version, gen id.Generator) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		p := path.Dir(info.FullMethod)[1:]
		if ts.IsObservable(p) {
			return handler(ctx, req)
		}

		md := ExtractIncoming(ctx)

		ctx = m.WithUserAgent(ctx, extractUserAgent(ctx, md, userAgent))

		id := extractRequestID(ctx, gen, md)
		ctx = m.WithRequestID(ctx, id)

		kind, ip := extractIPAddr(ctx, md)
		ctx = m.WithIPAddr(ctx, ip)
		ctx = m.WithIPAddrKind(ctx, kind)

		ctx = m.WithGeolocation(ctx, extractGeolocation(md))
		ctx = m.WithAuthorization(ctx, extractAuthorization(ctx, md))

		_ = grpc.SetHeader(ctx, metadata.Pairs("service-version", version.String(), "request-id", id.Value()))

		return handler(ctx, req)
	}
}

// StreamServerInterceptor for meta.
func StreamServerInterceptor(userAgent env.UserAgent, version env.Version, gen id.Generator) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		p := path.Dir(info.FullMethod)[1:]
		if ts.IsObservable(p) {
			return handler(srv, stream)
		}

		_ = stream.SetHeader(metadata.Pairs("service-version", version.String()))

		ctx := stream.Context()
		md := ExtractIncoming(ctx)
		ctx = m.WithUserAgent(ctx, extractUserAgent(ctx, md, userAgent))

		id := extractRequestID(ctx, gen, md)
		_ = stream.SetHeader(metadata.Pairs("request-id", id.Value()))

		ctx = m.WithRequestID(ctx, id)

		kind, ip := extractIPAddr(ctx, md)
		ctx = m.WithIPAddr(ctx, ip)
		ctx = m.WithIPAddrKind(ctx, kind)

		ctx = m.WithGeolocation(ctx, extractGeolocation(md))
		ctx = m.WithAuthorization(ctx, extractAuthorization(ctx, md))

		wrappedStream := middleware.WrapServerStream(stream)
		wrappedStream.WrappedContext = ctx

		return handler(srv, wrappedStream)
	}
}

// UnaryClientInterceptor for meta.
func UnaryClientInterceptor(userAgent env.UserAgent, gen id.Generator) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, fullMethod string, req, resp any, conn *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		md := ExtractOutgoing(ctx)

		ua := extractUserAgent(ctx, md, userAgent)
		ctx = m.WithUserAgent(ctx, ua)
		md.Append("user-agent", ua.Value())

		id := extractRequestID(ctx, gen, md)
		ctx = m.WithRequestID(ctx, id)
		md.Append("request-id", id.Value())

		ctx = metadata.NewOutgoingContext(ctx, md)

		return invoker(ctx, fullMethod, req, resp, conn, opts...)
	}
}

// StreamClientInterceptor for meta.
func StreamClientInterceptor(userAgent env.UserAgent, gen id.Generator) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, conn *grpc.ClientConn, fullMethod string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		md := ExtractOutgoing(ctx)

		ua := extractUserAgent(ctx, md, userAgent)
		ctx = m.WithUserAgent(ctx, ua)
		md.Append("user-agent", ua.Value())

		id := extractRequestID(ctx, gen, md)
		ctx = m.WithRequestID(ctx, id)
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
	if ua := m.UserAgent(ctx); ua.Value() != "" {
		return ua
	}

	if ua := md.Get("user-agent"); len(ua) > 0 {
		return meta.String(ua[0])
	}

	return meta.String(userAgent.String())
}

func extractRequestID(ctx context.Context, gen id.Generator, md metadata.MD) meta.Value {
	if id := m.RequestID(ctx); id.Value() != "" {
		return id
	}

	if id := md.Get("request-id"); len(id) > 0 {
		return meta.String(id[0])
	}

	return meta.String(gen.Generate())
}

func extractAuthorization(ctx context.Context, md metadata.MD) meta.Value {
	a := authorization(md)
	if a == "" {
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
