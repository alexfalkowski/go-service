package meta

import (
	"context"

	"github.com/alexfalkowski/go-service/pkg/transport/meta"
	"github.com/google/uuid"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// UnaryServerInterceptor for meta.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md := ExtractIncoming(ctx)

		userAgent := extractUserAgent(ctx, md)
		ctx = meta.WithUserAgent(ctx, userAgent)

		requestID := extractRequestID(ctx, md)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx = meta.WithRequestID(ctx, requestID)

		headers := metadata.Pairs("request-id", requestID, "ua", userAgent)
		if err := grpc.SendHeader(ctx, headers); err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

// StreamServerInterceptor for meta.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := stream.Context()
		md := ExtractIncoming(ctx)

		userAgent := extractUserAgent(ctx, md)
		ctx = meta.WithUserAgent(ctx, userAgent)

		requestID := extractRequestID(ctx, md)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx = meta.WithRequestID(ctx, requestID)

		headers := metadata.Pairs("request-id", requestID, "ua", userAgent)
		if err := grpc.SendHeader(ctx, headers); err != nil {
			return err
		}

		wrappedStream := grpcMiddleware.WrapServerStream(stream)
		wrappedStream.WrappedContext = ctx

		return handler(srv, stream)
	}
}

// UnaryClientInterceptor for meta.
func UnaryClientInterceptor(userAgent string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, fullMethod string, req, resp interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
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

		ctx = metadata.AppendToOutgoingContext(ctx, "request-id", requestID, "ua", ua)

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

		ctx = metadata.AppendToOutgoingContext(ctx, "request-id", requestID, "ua", ua)

		return streamer(ctx, desc, cc, fullMethod, opts...)
	}
}

func extractUserAgent(ctx context.Context, md metadata.MD) string {
	if mdUserAgent := md.Get("ua"); len(mdUserAgent) > 0 {
		return mdUserAgent[0]
	}

	return meta.UserAgent(ctx)
}

func extractRequestID(ctx context.Context, md metadata.MD) string {
	if mdRequestID := md.Get("request-id"); len(mdRequestID) > 0 {
		return mdRequestID[0]
	}

	return meta.RequestID(ctx)
}
