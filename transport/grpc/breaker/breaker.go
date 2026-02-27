package breaker

import (
	"github.com/alexfalkowski/go-service/v2/breaker"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/alexfalkowski/go-service/v2/sync"
)

// Settings is an alias for `github.com/alexfalkowski/go-service/v2/breaker.Settings`.
//
// It is re-exported from this package so callers can configure breaker behavior (trip thresholds, timeouts,
// half-open probing, etc.) without importing the lower-level breaker package directly.
type Settings = breaker.Settings

// UnaryClientInterceptor returns a gRPC unary client interceptor guarded by circuit breakers.
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
// codes as failures (see `WithFailureCodes` and the defaults in `defaultOpts`). Calls that return other codes
// do not contribute to opening the breaker.
//
// # Error mapping
//
// If the breaker rejects a call because it is open (`breaker.ErrOpenState`) or because the half-open
// MaxRequests limit would be exceeded (`breaker.ErrTooManyRequests`), the interceptor maps that condition to
// a gRPC `ResourceExhausted` status error.
//
// All other errors from the invoker are returned as-is.
func UnaryClientInterceptor(options ...Option) grpc.UnaryClientInterceptor {
	o := defaultOpts()
	for _, option := range options {
		option.apply(o)
	}

	r := &registry{opts: o, breakers: sync.NewMap[string, *breaker.CircuitBreaker]()}

	return func(ctx context.Context, fullMethod string, req, resp any, conn *grpc.ClientConn, invoker grpc.UnaryInvoker, callOpts ...grpc.CallOption) error {
		cb := r.get(fullMethod)
		_, err := cb.Execute(func() (any, error) {
			return nil, invoker(ctx, fullMethod, req, resp, conn, callOpts...)
		})
		if err != nil {
			if errors.Is(err, breaker.ErrOpenState) || errors.Is(err, breaker.ErrTooManyRequests) {
				return status.Error(codes.ResourceExhausted, err.Error())
			}

			return err
		}
		return nil
	}
}
