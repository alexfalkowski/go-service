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

// Settings is an alias for breaker.Settings.
type Settings = breaker.Settings

// UnaryClientInterceptor returns a gRPC unary client interceptor guarded by a circuit breaker.
//
// A separate circuit breaker is maintained per fullMethod.
// Whether an invocation is counted as successful is determined by gRPC status code classification
// (configurable via options; see default failure codes in this package).
//
// If the circuit breaker is open or has too many requests, the interceptor returns ResourceExhausted.
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
