package limiter

import (
	"context"
	"path"

	"github.com/alexfalkowski/go-service/limiter"
	"github.com/alexfalkowski/go-service/transport/strings"
	l "github.com/sethvargo/go-limiter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor for gRPC.
func UnaryServerInterceptor(limiter l.Store, key limiter.KeyFunc) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		p := path.Dir(info.FullMethod)[1:]
		if strings.IsObservable(p) {
			return handler(ctx, req)
		}

		if err := limit(ctx, limiter, key); err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func limit(ctx context.Context, store l.Store, key limiter.KeyFunc) error {
	ok, info, err := limiter.Take(ctx, store, key)
	if err != nil {
		return internalError(err)
	}

	if err := grpc.SetHeader(ctx, metadata.Pairs("ratelimit", info)); err != nil {
		return internalError(err)
	}

	if !ok {
		return status.Errorf(codes.ResourceExhausted, "limiter: resource exhausted, %s", info)
	}

	return nil
}

func internalError(err error) error {
	return status.Errorf(codes.Internal, "limiter: %s", err.Error())
}
