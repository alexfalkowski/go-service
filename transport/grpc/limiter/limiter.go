package limiter

import (
	"context"
	"fmt"
	"path"

	"github.com/alexfalkowski/go-service/limiter"
	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/transport/strings"
	l "github.com/sethvargo/go-limiter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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

func limit(ctx context.Context, limiter l.Store, key limiter.KeyFunc) error {
	tokens, remaining, reset, ok, err := limiter.Take(ctx, meta.ValueOrBlank(key(ctx)))
	if err != nil {
		return status.Errorf(codes.Internal, "limiter: %s", err.Error())
	}

	if !ok {
		return status.Errorf(codes.ResourceExhausted, fmt.Sprintf("limiter: limit=%d, remaining=%d, reset=%d", tokens, remaining, reset))
	}

	return nil
}
