package limiter

import (
	"context"

	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewLimiter for gRPC.
func NewLimiter(formatted string) (*limiter.Limiter, error) {
	rate, err := limiter.NewRateFromFormatted(formatted)
	if err != nil {
		return nil, err
	}

	store := memory.NewStore()

	return limiter.New(store, rate), nil
}

// KeyFunc to get for gRPC.
type KeyFunc func(context.Context) string

// UnaryServerInterceptor for gRPC.
func UnaryServerInterceptor(limiter *limiter.Limiter, key KeyFunc) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
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
func StreamServerInterceptor(limiter *limiter.Limiter, key KeyFunc) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
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
