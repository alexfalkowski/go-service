package limiter

import (
	"context"
	"path"

	"github.com/alexfalkowski/go-service/limiter"
	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/transport/strings"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	v3 "github.com/ulule/limiter/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor for gRPC.
func UnaryServerInterceptor(limiter *v3.Limiter, key limiter.KeyFunc) grpc.UnaryServerInterceptor {
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

// StreamServerInterceptor for gRPC.
func StreamServerInterceptor(limiter *v3.Limiter, key limiter.KeyFunc) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		p := path.Dir(info.FullMethod)[1:]
		if strings.IsObservable(p) {
			return handler(srv, stream)
		}

		ctx := stream.Context()
		if err := limit(ctx, limiter, key); err != nil {
			return err
		}

		wrappedStream := middleware.WrapServerStream(stream)
		wrappedStream.WrappedContext = ctx

		return handler(srv, stream)
	}
}

func limit(ctx context.Context, limiter *v3.Limiter, key limiter.KeyFunc) error {
	// Memory stores do not return error.
	context, _ := limiter.Get(ctx, meta.ValueOrBlank(key(ctx)))

	if context.Reached {
		return status.Errorf(codes.ResourceExhausted, "limit: %d allowed", context.Limit)
	}

	return nil
}
