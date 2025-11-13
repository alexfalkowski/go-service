package limiter

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/limiter"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/meta"
	"github.com/alexfalkowski/go-service/v2/transport/strings"
)

// KeyMap is just an alias for limiter.KeyMap.
type KeyMap = limiter.KeyMap

// NewServerLimiter for gRPC.
func NewServerLimiter(lc di.Lifecycle, keys KeyMap, cfg *limiter.Config) (*Server, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}

	limiter, err := limiter.NewLimiter(lc, keys, cfg)
	if err != nil {
		return nil, err
	}

	return &Server{limiter}, nil
}

// Server limiter.
type Server struct {
	*limiter.Limiter
}

// UnaryServerInterceptor for limiter.
func UnaryServerInterceptor(limiter *Server) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if strings.IsObservable(info.FullMethod) {
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

// NewClientLimiter for gRPC.
func NewClientLimiter(lc di.Lifecycle, keys KeyMap, cfg *limiter.Config) (*Client, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}

	limiter, err := limiter.NewLimiter(lc, keys, cfg)
	if err != nil {
		return nil, err
	}

	return &Client{limiter}, nil
}

// Client limiter.
type Client struct {
	*limiter.Limiter
}

// UnaryClientInterceptor for limiter.
func UnaryClientInterceptor(limiter *Client) grpc.UnaryClientInterceptor {
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
