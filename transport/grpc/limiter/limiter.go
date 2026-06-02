package limiter

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/net/grpc/meta"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/alexfalkowski/go-service/v2/net/grpc/strings"
	"github.com/alexfalkowski/go-service/v2/transport/limiter"
)

// KeyMap is an alias for [limiter.KeyMap].
//
// It maps limiter key kinds (for example, "user-agent" or "ip") to functions that derive a rate-limit key
// from the request context.
type KeyMap = limiter.KeyMap

// NewServerLimiter constructs a gRPC server-side rate limiter.
//
// If cfg is disabled, it returns (nil, nil) so callers can treat the limiter as not configured.
//
// The returned limiter is backed by [limiter.NewLimiter] and is registered with the provided lifecycle.
// The `keys` map controls how request contexts are turned into limiter keys (for example, per user-agent).
func NewServerLimiter(lc di.Lifecycle, keys KeyMap, cfg *limiter.Config) (*Server, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}

	rateLimiter, err := limiter.NewLimiter(lc, keys, cfg)
	if err != nil {
		return nil, err
	}

	return &Server{rateLimiter}, nil
}

// Server wraps *[limiter.Limiter] for gRPC server integration.
type Server struct {
	*limiter.Limiter
}

// UnaryServerInterceptor returns a gRPC unary server interceptor that enforces rate limiting.
//
// Operation unary RPC methods (health/metrics/etc.) bypass limiting (see [github.com/alexfalkowski/go-service/v2/net/grpc/strings.IsOperationMethod]).
// Stream RPCs do not bypass limiting because long-lived streams, such as health Watch, can hold server resources
// until the client disconnects.
//
// On every request, the interceptor calls `limiter.Take(ctx)` to determine whether the request is allowed:
//
//   - If `Take` returns an error, the interceptor returns [codes.Internal].
//   - If `Take` returns a header string, it is attached to response metadata as the "ratelimit" header.
//   - If the request is not allowed, the interceptor returns [codes.ResourceExhausted].
//   - Otherwise, it invokes the handler.
//
// Callers should only install this interceptor when limiter is non-nil.
func UnaryServerInterceptor(limiter *Server) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if strings.IsOperationMethod(info.FullMethod) {
			return handler(ctx, req)
		}

		ok, header, err := limiter.Take(ctx)
		if err != nil {
			return nil, status.SafeError(codes.Internal, err)
		}

		_ = grpc.SetHeader(ctx, meta.Pairs("ratelimit", header))

		if !ok {
			return nil, status.Error(codes.ResourceExhausted, grpc.StatusText(codes.ResourceExhausted))
		}

		return handler(ctx, req)
	}
}

// StreamServerInterceptor returns a gRPC stream server interceptor that enforces rate limiting.
//
// Unlike unary RPCs, operation streams are limited. This keeps long-lived streams, such as health Watch, from
// bypassing the resource controls that protect regular service streams.
//
// The limiter is admission-only for streams: it takes one token when the stream is opened and does not meter
// individual messages or stream lifetime. Use gRPC server options such as max_concurrent_streams, plus edge,
// gateway, ingress, load-balancer, or service-mesh limits when long-lived stream occupancy needs a hard cap.
//
// On every stream, the interceptor calls `limiter.Take(ctx)` to determine whether the stream is allowed:
//
//   - If `Take` returns an error, the interceptor returns [codes.Internal].
//   - If `Take` returns a header string, it is attached to response metadata as the "ratelimit" header.
//   - If the stream is not allowed, the interceptor returns [codes.ResourceExhausted].
//   - Otherwise, it invokes the handler.
//
// Callers should only install this interceptor when limiter is non-nil.
func StreamServerInterceptor(limiter *Server) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ok, header, err := limiter.Take(stream.Context())
		if err != nil {
			return status.SafeError(codes.Internal, err)
		}

		_ = stream.SetHeader(meta.Pairs("ratelimit", header))

		if !ok {
			return status.Error(codes.ResourceExhausted, grpc.StatusText(codes.ResourceExhausted))
		}

		return handler(srv, stream)
	}
}

// NewClientLimiter constructs a gRPC client-side rate limiter.
//
// If cfg is disabled, it returns (nil, nil) so callers can treat the limiter as not configured.
//
// The returned limiter is backed by [limiter.NewLimiter] and is registered with the provided lifecycle.
// The `keys` map controls how request contexts are turned into limiter keys.
func NewClientLimiter(lc di.Lifecycle, keys KeyMap, cfg *limiter.Config) (*Client, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}

	rateLimiter, err := limiter.NewLimiter(lc, keys, cfg)
	if err != nil {
		return nil, err
	}

	return &Client{rateLimiter}, nil
}

// Client wraps *[limiter.Limiter] for gRPC client integration.
type Client struct {
	*limiter.Limiter
}

// UnaryClientInterceptor returns a gRPC unary client interceptor that enforces rate limiting.
//
// The interceptor calls `limiter.Take(ctx)` before invoking the RPC:
//
//   - If `Take` returns an error, it returns [codes.Internal].
//   - If the request is not allowed, it returns [codes.ResourceExhausted].
//   - Otherwise, it invokes the underlying `invoker`.
//
// Callers should only install this interceptor when limiter is non-nil.
func UnaryClientInterceptor(limiter *Client) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, fullMethod string, req, resp any, conn *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ok, _, err := limiter.Take(ctx)
		if err != nil {
			return status.SafeError(codes.Internal, err)
		}

		if !ok {
			return status.Error(codes.ResourceExhausted, grpc.StatusText(codes.ResourceExhausted))
		}

		return invoker(ctx, fullMethod, req, resp, conn, opts...)
	}
}

// StreamClientInterceptor returns a gRPC stream client interceptor that enforces rate limiting.
//
// The interceptor calls `limiter.Take(ctx)` before opening the stream:
//
//   - If `Take` returns an error, it returns [codes.Internal].
//   - If the stream is not allowed, it returns [codes.ResourceExhausted].
//   - Otherwise, it invokes the underlying `streamer`.
//
// Callers should only install this interceptor when limiter is non-nil.
func StreamClientInterceptor(limiter *Client) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, conn *grpc.ClientConn, fullMethod string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		ok, _, err := limiter.Take(ctx)
		if err != nil {
			return nil, status.SafeError(codes.Internal, err)
		}

		if !ok {
			return nil, status.Error(codes.ResourceExhausted, grpc.StatusText(codes.ResourceExhausted))
		}

		return streamer(ctx, desc, conn, fullMethod, opts...)
	}
}
