package meta

import (
	"context"
	"fmt"

	"github.com/alexfalkowski/go-service/meta"
	m "github.com/alexfalkowski/go-service/transport/meta"
	"github.com/google/uuid"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// UnaryServerInterceptor for meta.
func UnaryServerInterceptor(userAgent string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md := ExtractIncoming(ctx)

		userAgent := extractUserAgent(ctx, md, userAgent)
		ctx = m.WithUserAgent(ctx, userAgent)

		requestID := extractRequestID(ctx, md)
		if meta.IsBlank(requestID) {
			requestID = uuid.New()
		}

		ctx = m.WithRequestID(ctx, requestID)

		return handler(ctx, req)
	}
}

// StreamServerInterceptor for meta.
func StreamServerInterceptor(userAgent string) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := stream.Context()
		md := ExtractIncoming(ctx)

		userAgent := extractUserAgent(ctx, md, userAgent)
		ctx = m.WithUserAgent(ctx, userAgent)

		requestID := extractRequestID(ctx, md)
		if meta.IsBlank(requestID) {
			requestID = uuid.New()
		}

		ctx = m.WithRequestID(ctx, requestID)

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
		if meta.IsBlank(ua) {
			ua = meta.String(userAgent)
		}

		ctx = m.WithUserAgent(ctx, ua)

		requestID := extractRequestID(ctx, md)
		if meta.IsBlank(requestID) {
			requestID = uuid.New()
		}

		ctx = m.WithRequestID(ctx, requestID)

		return invoker(ctx, fullMethod, req, resp, cc, opts...)
	}
}

// StreamClientInterceptor for meta.
func StreamClientInterceptor(userAgent string) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, fullMethod string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		md := ExtractOutgoing(ctx)

		ua := extractUserAgent(ctx, md, userAgent)
		if meta.IsBlank(ua) {
			ua = meta.String(userAgent)
		}

		ctx = m.WithUserAgent(ctx, ua)

		requestID := extractRequestID(ctx, md)
		if meta.IsBlank(requestID) {
			requestID = uuid.New()
		}

		ctx = m.WithRequestID(ctx, requestID)

		return streamer(ctx, desc, cc, fullMethod, opts...)
	}
}

func extractUserAgent(ctx context.Context, md metadata.MD, userAgent string) fmt.Stringer {
	if ua := md.Get(runtime.MetadataPrefix + "user-agent"); len(ua) > 0 {
		return meta.String(ua[0])
	}

	if ua := md.Get("user-agent"); len(ua) > 0 {
		return meta.String(ua[0])
	}

	if u := m.UserAgent(ctx); u != nil {
		return u
	}

	return meta.String(userAgent)
}

func extractRequestID(ctx context.Context, md metadata.MD) fmt.Stringer {
	if id := md.Get("request-id"); len(id) > 0 {
		return meta.String(id[0])
	}

	return m.RequestID(ctx)
}
