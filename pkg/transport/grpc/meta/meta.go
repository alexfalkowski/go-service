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

		var userAgent string

		if mdUserAgent := md.Get("user-agent"); len(mdUserAgent) > 0 {
			userAgent = mdUserAgent[0]
		}

		ctx = meta.WithUserAgent(ctx, userAgent)

		var requestID string

		if mdRequestID := md.Get("request-id"); len(mdRequestID) > 0 {
			requestID = mdRequestID[0]
		}

		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx = meta.WithRequestID(ctx, requestID)
		header := metadata.Pairs("request-id", requestID)

		if err := grpc.SendHeader(ctx, header); err != nil {
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

		var userAgent string

		if mdUserAgent := md.Get("user-agent"); len(mdUserAgent) > 0 {
			userAgent = mdUserAgent[0]
		}

		ctx = meta.WithUserAgent(ctx, userAgent)

		var requestID string

		if mdRequestID := md.Get("request-id"); len(mdRequestID) > 0 {
			requestID = mdRequestID[0]
		}

		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx = meta.WithRequestID(ctx, requestID)
		header := metadata.Pairs("request-id", requestID)

		if err := grpc.SendHeader(ctx, header); err != nil {
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
		ctx = meta.WithUserAgent(ctx, userAgent)

		requestID := meta.RequestID(ctx)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx = meta.WithRequestID(ctx, requestID)
		ctx = metadata.AppendToOutgoingContext(ctx, "request-id", requestID)

		return invoker(ctx, fullMethod, req, resp, cc, opts...)
	}
}

// StreamClientInterceptor for meta.
func StreamClientInterceptor(userAgent string) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, fullMethod string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		ctx = meta.WithUserAgent(ctx, userAgent)

		requestID := meta.RequestID(ctx)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx = meta.WithRequestID(ctx, requestID)
		ctx = metadata.AppendToOutgoingContext(ctx, "request-id", requestID)

		return streamer(ctx, desc, cc, fullMethod, opts...)
	}
}
