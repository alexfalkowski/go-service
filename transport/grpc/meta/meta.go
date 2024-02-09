package meta

import (
	"context"

	"github.com/alexfalkowski/go-service/transport/meta"
	"github.com/google/uuid"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// UnaryServerInterceptor for meta.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md := ExtractIncoming(ctx)

		userAgent := extractUserAgent(ctx, md)
		ctx = meta.WithUserAgent(ctx, userAgent)

		requestID := extractRequestID(ctx, md)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx = meta.WithRequestID(ctx, requestID)

		return handler(ctx, req)
	}
}

// StreamServerInterceptor for meta.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := stream.Context()
		md := ExtractIncoming(ctx)

		userAgent := extractUserAgent(ctx, md)
		ctx = meta.WithUserAgent(ctx, userAgent)

		requestID := extractRequestID(ctx, md)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx = meta.WithRequestID(ctx, requestID)

		wrappedStream := middleware.WrapServerStream(stream)
		wrappedStream.WrappedContext = ctx

		return handler(srv, wrappedStream)
	}
}

// UnaryClientInterceptor for meta.
func UnaryClientInterceptor(userAgent string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, fullMethod string, req, resp any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		md := ExtractOutgoing(ctx)

		ua := extractUserAgent(ctx, md)
		if ua == "" {
			ua = userAgent
		}

		ctx = meta.WithUserAgent(ctx, ua)

		requestID := extractRequestID(ctx, md)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx = meta.WithRequestID(ctx, requestID)

		return invoker(ctx, fullMethod, req, resp, cc, opts...)
	}
}

// StreamClientInterceptor for meta.
func StreamClientInterceptor(userAgent string) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, fullMethod string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		md := ExtractOutgoing(ctx)

		ua := extractUserAgent(ctx, md)
		if ua == "" {
			ua = userAgent
		}

		ctx = meta.WithUserAgent(ctx, ua)

		requestID := extractRequestID(ctx, md)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx = meta.WithRequestID(ctx, requestID)

		return streamer(ctx, desc, cc, fullMethod, opts...)
	}
}

func extractUserAgent(ctx context.Context, md metadata.MD) string {
	if ua := md.Get(runtime.MetadataPrefix + "user-agent"); len(ua) > 0 {
		return ua[0]
	}

	if ua := md.Get("user-agent"); len(ua) > 0 {
		return ua[0]
	}

	return meta.UserAgent(ctx)
}

func extractRequestID(ctx context.Context, md metadata.MD) string {
	if id := md.Get("request-id"); len(id) > 0 {
		return id[0]
	}

	return meta.RequestID(ctx)
}
