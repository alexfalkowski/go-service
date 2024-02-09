package limiter

import (
	"context"

	"github.com/alexfalkowski/go-service/limiter"
	l "github.com/ulule/limiter/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor for gRPC.
func UnaryServerInterceptor(limiter *l.Limiter, key limiter.KeyFunc) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		context, err := limiter.Get(ctx, key(ctx))
		if err != nil {
			return nil, err
		}

		if context.Reached {
			return nil, status.Errorf(codes.ResourceExhausted, "limit: %d allowed", context.Limit)
		}

		return handler(ctx, req)
	}
}

// StreamServerInterceptor for gRPC.
func StreamServerInterceptor(limiter *l.Limiter, key limiter.KeyFunc) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := stream.Context()

		context, err := limiter.Get(ctx, key(ctx))
		if err != nil {
			return err
		}

		if context.Reached {
			return status.Errorf(codes.ResourceExhausted, "limit: %d allowed", context.Limit)
		}

		return handler(srv, stream)
	}
}
