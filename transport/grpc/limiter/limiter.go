package limiter

import (
	"context"
	"path"

	"github.com/alexfalkowski/go-service/v2/limiter"
	"github.com/alexfalkowski/go-service/v2/transport/strings"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Limiter is just an alias for limiter.Limiter.
type Limiter = limiter.Limiter

// UnaryServerInterceptor for gRPC.
func UnaryServerInterceptor(limiter *Limiter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		p := path.Dir(info.FullMethod)[1:]
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
		return internalError(err)
	}

	_ = grpc.SetHeader(ctx, metadata.Pairs("ratelimit", info))

	if !ok {
		return status.Errorf(codes.ResourceExhausted, "limiter: resource exhausted, %s", info)
	}

	return nil
}

func internalError(err error) error {
	return status.Errorf(codes.Internal, "limiter: %s", err.Error())
}
