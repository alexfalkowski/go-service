package limiter

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/limiter"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/meta"
	"github.com/alexfalkowski/go-service/v2/transport/strings"
)

// Limiter is just an alias for limiter.Limiter.
type Limiter = limiter.Limiter

// UnaryServerInterceptor for limiter.
func UnaryServerInterceptor(limiter *Limiter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		p := info.FullMethod[1:]
		if strings.IsObservable(p) {
			return handler(ctx, req)
		}

		ok, header, err := limiter.Take(ctx)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "limiter: %s", err.Error())
		}

		_ = grpc.SetHeader(ctx, meta.Pairs("ratelimit", header))

		if !ok {
			return nil, status.Errorf(codes.ResourceExhausted, "limiter: resource exhausted, %s", header)
		}

		return handler(ctx, req)
	}
}

// UnaryClientInterceptor for limiter.
func UnaryClientInterceptor(limiter *Limiter) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, fullMethod string, req, resp any, conn *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ok, header, err := limiter.Take(ctx)
		if err != nil {
			return status.Errorf(codes.Internal, "limiter: %s", err.Error())
		}

		if !ok {
			return status.Errorf(codes.ResourceExhausted, "limiter: resource exhausted, %s", header)
		}

		return invoker(ctx, fullMethod, req, resp, conn, opts...)
	}
}
