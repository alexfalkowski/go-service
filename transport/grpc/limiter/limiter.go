package limiter

import (
	"context"

	"github.com/alexfalkowski/go-service/v2/limiter"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/meta"
	"github.com/alexfalkowski/go-service/v2/transport/strings"
)

// Limiter is just an alias for limiter.Limiter.
type Limiter = limiter.Limiter

// UnaryServerInterceptor for gRPC.
func UnaryServerInterceptor(limiter *Limiter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		p := info.FullMethod[1:]
		if strings.IsObservable(p) {
			return handler(ctx, req)
		}

		if err := limit(ctx, limiter); err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func limit(ctx context.Context, limiter *Limiter) error {
	ok, info, err := limiter.Take(ctx)
	if err != nil {
		return status.Errorf(codes.Internal, "limiter: %s", err.Error())
	}

	_ = grpc.SetHeader(ctx, meta.Pairs("ratelimit", info))

	if !ok {
		return status.Errorf(codes.ResourceExhausted, "limiter: resource exhausted, %s", info)
	}

	return nil
}
