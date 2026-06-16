package limiter

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/transport/limiter"
)

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
// The interceptor calls `limiter.TakeDecision(ctx)` before invoking the RPC:
//
//   - If `TakeDecision` returns an error, it returns [codes.Internal].
//   - If the request is not allowed, it returns [codes.ResourceExhausted].
//   - Otherwise, it invokes the underlying `invoker`.
//
// Callers should only install this interceptor when limiter is non-nil.
func UnaryClientInterceptor(limiter *Client) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, fullMethod string, req, resp any, conn *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		decision, err := take(ctx, limiter.Limiter)
		if err != nil {
			return err
		}

		if !decision.Allowed() {
			return limitError()
		}

		return invoker(ctx, fullMethod, req, resp, conn, opts...)
	}
}

// StreamClientInterceptor returns a gRPC stream client interceptor that enforces rate limiting.
//
// The interceptor calls `limiter.TakeDecision(ctx)` before opening the stream, then meters each sent and received
// stream message:
//
//   - If `TakeDecision` returns an error, it returns [codes.Internal].
//   - If the stream is not allowed, it returns [codes.ResourceExhausted].
//   - Otherwise, it invokes the underlying `streamer`.
//
// Callers should only install this interceptor when limiter is non-nil.
func StreamClientInterceptor(limiter *Client) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, conn *grpc.ClientConn, fullMethod string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		decision, err := take(ctx, limiter.Limiter)
		if err != nil {
			return nil, err
		}

		if !decision.Allowed() {
			return nil, limitError()
		}

		stream, err := streamer(ctx, desc, conn, fullMethod, opts...)
		if err != nil {
			return nil, err
		}

		return &clientStream{ClientStream: stream, ctx: ctx, limiter: limiter.Limiter}, nil
	}
}

type clientStream struct {
	grpc.ClientStream
	ctx     context.Context
	limiter *limiter.Limiter
}

func (s *clientStream) RecvMsg(m any) error {
	if err := s.ClientStream.RecvMsg(m); err != nil {
		return err
	}

	decision, err := take(s.ctx, s.limiter)
	if err != nil {
		return err
	}

	if !decision.Allowed() {
		return limitError()
	}

	return nil
}

func (s *clientStream) SendMsg(m any) error {
	decision, err := take(s.ctx, s.limiter)
	if err != nil {
		return err
	}

	if !decision.Allowed() {
		return limitError()
	}

	return s.ClientStream.SendMsg(m)
}
