package meta

import (
	"context"
	"net"
	"path"
	"strings"

	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/security/header"
	m "github.com/alexfalkowski/go-service/transport/meta"
	ts "github.com/alexfalkowski/go-service/transport/strings"
	"github.com/google/uuid"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// UnaryServerInterceptor for meta.
func UnaryServerInterceptor(userAgent string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		p := path.Dir(info.FullMethod)[1:]
		if ts.IsObservable(p) {
			return handler(ctx, req)
		}

		md := ExtractIncoming(ctx)

		ctx = m.WithUserAgent(ctx, extractUserAgent(ctx, md, userAgent))
		ctx = m.WithRequestID(ctx, extractRequestID(ctx, md))

		kind, ip := extractIPAddr(ctx, md)
		ctx = m.WithIPAddr(ctx, ip)
		ctx = m.WithIPAddrKind(ctx, kind)

		ctx = m.WithGeolocation(ctx, extractGeolocation(ctx, md))
		ctx = m.WithAuthorization(ctx, extractAuthorization(ctx, md))

		return handler(ctx, req)
	}
}

// StreamServerInterceptor for meta.
func StreamServerInterceptor(userAgent string) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		p := path.Dir(info.FullMethod)[1:]
		if ts.IsObservable(p) {
			return handler(srv, stream)
		}

		ctx := stream.Context()
		md := ExtractIncoming(ctx)

		ctx = m.WithUserAgent(ctx, extractUserAgent(ctx, md, userAgent))
		ctx = m.WithRequestID(ctx, extractRequestID(ctx, md))

		kind, ip := extractIPAddr(ctx, md)
		ctx = m.WithIPAddr(ctx, ip)
		ctx = m.WithIPAddrKind(ctx, kind)

		ctx = m.WithGeolocation(ctx, extractGeolocation(ctx, md))
		ctx = m.WithAuthorization(ctx, extractAuthorization(ctx, md))

		wrappedStream := middleware.WrapServerStream(stream)
		wrappedStream.WrappedContext = ctx

		return handler(srv, wrappedStream)
	}
}

// UnaryClientInterceptor for meta.
func UnaryClientInterceptor(userAgent string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, fullMethod string, req, resp any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		md := ExtractOutgoing(ctx)

		ua := extractUserAgent(ctx, md, userAgent)
		ctx = m.WithUserAgent(ctx, ua)
		md.Append("user-agent", ua.Value())

		id := extractRequestID(ctx, md)
		ctx = m.WithRequestID(ctx, extractRequestID(ctx, md))
		md.Append("request-id", id.Value())

		ctx = metadata.NewOutgoingContext(ctx, md)

		return invoker(ctx, fullMethod, req, resp, cc, opts...)
	}
}

// StreamClientInterceptor for meta.
func StreamClientInterceptor(userAgent string) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, fullMethod string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		md := ExtractOutgoing(ctx)

		ua := extractUserAgent(ctx, md, userAgent)
		ctx = m.WithUserAgent(ctx, ua)
		md.Append("user-agent", ua.Value())

		id := extractRequestID(ctx, md)
		ctx = m.WithRequestID(ctx, extractRequestID(ctx, md))
		md.Append("request-id", id.Value())

		ctx = metadata.NewOutgoingContext(ctx, md)

		return streamer(ctx, desc, cc, fullMethod, opts...)
	}
}

func extractIPAddr(ctx context.Context, md metadata.MD) (meta.Valuer, meta.Valuer) {
	headers := []string{"x-real-ip", "cf-connecting-ip", "true-client-ip", "x-forwarded-for"}
	for _, k := range headers {
		if f := md.Get(k); len(f) > 0 {
			return meta.String(k), meta.String(strings.Split(f[0], ",")[0])
		}
	}

	peerKind := meta.String("peer")

	p, ok := peer.FromContext(ctx)
	if !ok {
		return peerKind, meta.Blank()
	}

	addr := p.Addr.String()

	host, _, err := net.SplitHostPort(p.Addr.String())
	if err != nil {
		return peerKind, meta.String(addr)
	}

	return peerKind, meta.String(host)
}

func extractUserAgent(ctx context.Context, md metadata.MD, userAgent string) meta.Valuer {
	if ua := m.UserAgent(ctx); ua != nil {
		return ua
	}

	if ua := md.Get("user-agent"); len(ua) > 0 {
		return meta.String(ua[0])
	}

	return meta.String(userAgent)
}

func extractRequestID(ctx context.Context, md metadata.MD) meta.Valuer {
	if id := m.RequestID(ctx); id != nil {
		return id
	}

	if id := md.Get("request-id"); len(id) > 0 {
		return meta.String(id[0])
	}

	return meta.ToString(uuid.New())
}

func extractAuthorization(ctx context.Context, md metadata.MD) meta.Valuer {
	a := authorization(md)
	if a == "" {
		return meta.Blank()
	}

	_, t, err := header.ParseAuthorization(a)
	if err != nil {
		meta.WithAttribute(ctx, "authError", meta.Error(err))

		return meta.Blank()
	}

	return meta.Ignored(t)
}

func authorization(md metadata.MD) string {
	if a := md.Get("authorization"); len(a) > 0 {
		return a[0]
	}

	return ""
}

func extractGeolocation(ctx context.Context, md metadata.MD) meta.Valuer {
	if gl := m.Geolocation(ctx); gl != nil {
		return gl
	}

	if id := md.Get("geolocation"); len(id) > 0 {
		return meta.String(id[0])
	}

	return meta.Blank()
}
