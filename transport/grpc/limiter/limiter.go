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

// KeyMap is an alias for `limiter.KeyMap`.
//
// It maps limiter key kinds (for example, "user-agent" or "ip") to functions that derive a rate-limit key
// from the request context.
type KeyMap = limiter.KeyMap

// NewServerLimiter constructs a gRPC server-side rate limiter.
//
// If cfg is disabled, it returns (nil, nil) so callers can treat the limiter as not configured.
//
// The returned limiter is backed by `limiter.NewLimiter` and is registered with the provided lifecycle.
// The `keys` map controls how request contexts are turned into limiter keys (for example, per user-agent).
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

// Server wraps `*limiter.Limiter` for gRPC server integration.
type Server struct {
	*limiter.Limiter
}

// UnaryServerInterceptor returns a gRPC unary server interceptor that enforces rate limiting.
//
// Ignorable RPC methods (health/metrics/etc.) bypass limiting (see `transport/strings.IsIgnorable`).
//
// On every request, the interceptor calls `limiter.Take(ctx)` to determine whether the request is allowed:
//
//   - If `Take` returns an error, the interceptor returns `codes.Internal`.
//   - If `Take` returns a header string, it is attached to response metadata as the "ratelimit" header.
//   - If the request is not allowed, the interceptor returns `codes.ResourceExhausted`.
//   - Otherwise, it invokes the handler.
//
// Callers should only install this interceptor when limiter is non-nil.
func UnaryServerInterceptor(limiter *Server) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if strings.IsIgnorable(info.FullMethod) {
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

// NewClientLimiter constructs a gRPC client-side rate limiter.
//
// If cfg is disabled, it returns (nil, nil) so callers can treat the limiter as not configured.
//
// The returned limiter is backed by `limiter.NewLimiter` and is registered with the provided lifecycle.
// The `keys` map controls how request contexts are turned into limiter keys.
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

// Client wraps `*limiter.Limiter` for gRPC client integration.
type Client struct {
	*limiter.Limiter
}

// UnaryClientInterceptor returns a gRPC unary client interceptor that enforces rate limiting.
//
// The interceptor calls `limiter.Take(ctx)` before invoking the RPC:
//
//   - If `Take` returns an error, it returns `codes.Internal`.
//   - If the request is not allowed, it returns `codes.ResourceExhausted`.
//   - Otherwise, it invokes the underlying `invoker`.
//
// Callers should only install this interceptor when limiter is non-nil.
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
