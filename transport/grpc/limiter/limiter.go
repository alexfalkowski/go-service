package limiter

import (
	"context"

	"github.com/alexfalkowski/go-service/limiter"
	"github.com/alexfalkowski/go-service/meta"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	l "github.com/ulule/limiter/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor for gRPC.
func UnaryServerInterceptor(limiter *l.Limiter, key limiter.KeyFunc) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if err := limit(ctx, limiter, key); err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

// StreamServerInterceptor for gRPC.
func StreamServerInterceptor(limiter *l.Limiter, key limiter.KeyFunc) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := stream.Context()
		if err := limit(ctx, limiter, key); err != nil {
			return err
		}

		wrappedStream := middleware.WrapServerStream(stream)
		wrappedStream.WrappedContext = ctx

		return handler(srv, stream)
	}
}

func limit(ctx context.Context, limiter *l.Limiter, key limiter.KeyFunc) error {
	// Memory stores do not return error.
	context, _ := limiter.Get(ctx, meta.ValueOrBlank(key(ctx)))

	if context.Reached {
		return status.Errorf(codes.ResourceExhausted, "limit: %d allowed", context.Limit)
	}

	return nil
}
