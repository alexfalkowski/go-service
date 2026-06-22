package limiter

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/meta"
	"github.com/alexfalkowski/go-service/v2/net/grpc/strings"
	"github.com/alexfalkowski/go-service/v2/transport/limiter"
)

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
// On every request, the interceptor calls `limiter.TakeDecision(ctx)` to determine whether the request is allowed:
//
//   - If `Take` returns an error, the interceptor returns [codes.Internal].
//   - It attaches "ratelimit" and "ratelimit-policy" response metadata describing the current decision.
//   - If the request is not allowed, the interceptor returns [codes.ResourceExhausted].
//   - Otherwise, it invokes the handler.
//
// Callers should only install this interceptor when limiter is non-nil.
func UnaryServerInterceptor(limiter *Server) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if strings.IsOperationMethod(info.FullMethod) {
			return handler(ctx, req)
		}

		decision, err := take(ctx, limiter.Limiter)
		if err != nil {
			return nil, err
		}

		setHeader(ctx, decision)
		if !decision.Allowed() {
			return nil, limitError()
		}

		return handler(ctx, req)
	}
}

// StreamServerInterceptor returns a gRPC stream server interceptor that enforces rate limiting.
//
// Unlike unary RPCs, operation streams are limited. This keeps long-lived streams, such as health Watch, from
// bypassing the resource controls that protect regular service streams.
//
// The limiter takes one token when the stream is opened, then meters each received and sent stream message.
// Use gRPC server options such as max_concurrent_streams, plus edge, gateway, ingress, load-balancer, or
// service-mesh limits when long-lived stream occupancy needs a hard cap.
//
// On every stream, the interceptor calls `limiter.TakeDecision(ctx)` to determine whether the stream is allowed:
//
//   - If `Take` returns an error, the interceptor returns [codes.Internal].
//   - It attaches "ratelimit" and "ratelimit-policy" response metadata describing the current decision.
//   - If the stream is not allowed, the interceptor returns [codes.ResourceExhausted].
//   - Otherwise, it invokes the handler.
//
// Callers should only install this interceptor when limiter is non-nil.
func StreamServerInterceptor(limiter *Server) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		decision, err := take(stream.Context(), limiter.Limiter)
		if err != nil {
			return err
		}

		setStreamHeader(stream, decision)
		if !decision.Allowed() {
			return limitError()
		}

		return handler(srv, &serverStream{ServerStream: stream, limiter: limiter.Limiter})
	}
}

type serverStream struct {
	grpc.ServerStream
	limiter *limiter.Limiter
}

func (s *serverStream) RecvMsg(m any) error {
	if err := s.ServerStream.RecvMsg(m); err != nil {
		return err
	}

	decision, err := take(s.Context(), s.limiter)
	if err != nil {
		return err
	}

	setStreamHeader(s, decision)
	if !decision.Allowed() {
		return limitError()
	}

	return nil
}

func (s *serverStream) SendMsg(m any) error {
	decision, err := take(s.Context(), s.limiter)
	if err != nil {
		return err
	}

	setStreamHeader(s, decision)
	if !decision.Allowed() {
		return limitError()
	}

	return s.ServerStream.SendMsg(m)
}

func setHeader(ctx context.Context, decision limiter.Decision) {
	_ = grpc.SetHeader(ctx, limiterMetadata(decision))
}

func setStreamHeader(stream grpc.ServerStream, decision limiter.Decision) {
	md := limiterMetadata(decision)
	// Streaming headers may already be sent after the first response message; use trailers for later
	// limiter decisions so clients can still observe the current quota state on stream errors.
	if err := stream.SetHeader(md); err != nil {
		stream.SetTrailer(md)
	}
}

func limiterMetadata(decision limiter.Decision) meta.Map {
	return meta.Pairs("ratelimit", decision.Header(), "ratelimit-policy", decision.PolicyHeader())
}
