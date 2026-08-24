package breaker

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/alexfalkowski/go-service/v2/transport/breaker"
	"github.com/alexfalkowski/go-sync"
)

// Settings is an alias for [github.com/alexfalkowski/go-service/v2/transport/breaker.Settings].
//
// It is re-exported from this package so callers can configure breaker behavior (trip thresholds, timeouts,
// half-open probing, etc.) without importing the lower-level breaker package directly.
type Settings = breaker.Settings

// NewClient returns a gRPC client guarded by circuit breakers.
func NewClient(opts ...Option) *Client {
	options := defaultOptions()
	for _, opt := range opts {
		opt.apply(options)
	}

	return &Client{registry: &registry{options: options, breakers: sync.NewMap[string, *breaker.CircuitBreaker]()}}
}

// Client provides gRPC client circuit-breaking interceptors.
type Client struct {
	registry *registry
}

// UnaryInterceptor returns a gRPC unary client interceptor guarded by circuit breakers.
//
// The interceptor wraps the outgoing unary invocation (`invoker`) in a circuit breaker execution.
// When the breaker is closed, calls flow through normally. When the breaker transitions open, new calls
// are rejected until the breaker allows half-open probing per its settings.
//
// # Breaker scope
//
// A separate circuit breaker is maintained per `fullMethod`, so each downstream RPC method is isolated.
// Breaker instances are created lazily on first use and then reused for subsequent calls to the same method.
//
// # Failure classification
//
// The interceptor counts failures based on gRPC status codes. By default it treats a subset of transient/server
// codes as failures (see [WithFailureCodes] and the defaults in `defaultOptions`). Calls that return other codes
// do not contribute to opening the breaker. A context that is already deadline-exceeded before the invocation
// bypasses breaker accounting, because it provides no downstream-health signal.
//
// # Error mapping
//
// If the breaker rejects a call because it is open ([breaker.ErrOpenState]) or because the half-open
// MaxRequests limit would be exceeded ([breaker.ErrTooManyRequests]), the interceptor maps that condition to
// a locally marked gRPC `ResourceExhausted` status error so retry middleware returns it terminally.
//
// All other errors from the invoker are returned as-is.
func (c *Client) UnaryInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, fullMethod string, req, resp any, conn *grpc.ClientConn, invoker grpc.UnaryInvoker, callOpts ...grpc.CallOption) error {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return invoker(ctx, fullMethod, req, resp, conn, callOpts...)
		}

		cb := c.registry.get(fullMethod)
		_, err := cb.Execute(func() (any, error) {
			return nil, invoker(ctx, fullMethod, req, resp, conn, callOpts...)
		})
		if err != nil {
			if errors.Is(err, breaker.ErrOpenState) || errors.Is(err, breaker.ErrTooManyRequests) {
				return status.LocalError(status.SafeError(codes.ResourceExhausted, err))
			}

			return err
		}
		return nil
	}
}
