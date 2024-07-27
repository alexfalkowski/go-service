package limiter

import (
	"context"
	"fmt"
	"path"
	"time"

	"github.com/alexfalkowski/go-service/limiter"
	"github.com/alexfalkowski/go-service/meta"
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

func limit(ctx context.Context, limiter l.Store, key limiter.KeyFunc) error {
	tokens, remaining, reset, ok, err := limiter.Take(ctx, meta.ValueOrBlank(key(ctx)))
	if err != nil {
		return status.Errorf(codes.Internal, "limiter: %s", err.Error())
	}

	r := time.Until(time.Unix(0, int64(reset)))
	v := fmt.Sprintf("limit=%d, remaining=%d, reset=%s", tokens, remaining, r)

	if err := grpc.SetHeader(ctx, metadata.Pairs("ratelimit", v)); err != nil {
		return err
	}

	if !ok {
		return status.Errorf(codes.ResourceExhausted, "limiter: %s", v)
	}

	return nil
}
