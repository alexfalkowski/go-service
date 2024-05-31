package meta

import (
	"context"
	"net"
	"strings"

	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/security/header"
	m "github.com/alexfalkowski/go-service/transport/meta"
	"github.com/google/uuid"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// IPKeys to get the IP of the caller.
var IPKeys = []string{"x-real-ip", "cf-connecting-ip", "true-client-ip", "x-forwarded-for"}

// UnaryServerInterceptor for meta.
func UnaryServerInterceptor(userAgent string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md := ExtractIncoming(ctx)

		ctx = m.WithUserAgent(ctx, extractUserAgent(ctx, md, userAgent))
		ctx = m.WithRequestID(ctx, extractRequestID(ctx, md))
		ctx = m.WithIPAddr(ctx, meta.Ignored(IPAddr(ctx, md)))
		ctx = m.WithAuthorization(ctx, extractAuthorization(ctx, md))

		return handler(ctx, req)
	}
}

// StreamServerInterceptor for meta.
func StreamServerInterceptor(userAgent string) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := stream.Context()
		md := ExtractIncoming(ctx)

		ctx = m.WithUserAgent(ctx, extractUserAgent(ctx, md, userAgent))
		ctx = m.WithRequestID(ctx, extractRequestID(ctx, md))
		ctx = m.WithIPAddr(ctx, meta.Ignored(IPAddr(ctx, md)))
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

// IPAddr for meta.
func IPAddr(ctx context.Context, md metadata.MD) string {
	for _, k := range IPKeys {
		if f := md.Get(k); len(f) > 0 {
			return strings.Split(f[0], ",")[0]
		}
	}

	p, ok := peer.FromContext(ctx)
	if !ok {
		return ""
	}

	addr := p.Addr.String()

	host, _, err := net.SplitHostPort(p.Addr.String())
	if err != nil {
		return addr
	}

	return host
}

func extractUserAgent(ctx context.Context, md metadata.MD, userAgent string) meta.Valuer {
	if ua := m.UserAgent(ctx); ua != nil {
		return ua
	}

	if ua := md.Get(runtime.MetadataPrefix + "user-agent"); len(ua) > 0 {
		return meta.String(ua[0])
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
	if a := md.Get(runtime.MetadataPrefix + "authorization"); len(a) > 0 {
		return a[0]
	}

	if a := md.Get("authorization"); len(a) > 0 {
		return a[0]
	}

	return ""
}
