package breaker

import (
	"sync"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/alexfalkowski/go-service/v2/time"
	breaker "github.com/sony/gobreaker"
)

// Settings is an alias for the breaker.Settings.
type Settings = breaker.Settings

// Option interface for configuring the circuit breaker.
type Option interface {
	apply(opts *opts)
}

type opts struct {
	settings     Settings
	failureCodes map[codes.Code]struct{}
}

type optionFunc func(*opts)

func (f optionFunc) apply(o *opts) {
	f(o)
}

// WithSettings sets the settings for the circuit breaker.
func WithSettings(s Settings) Option {
	return optionFunc(func(o *opts) {
		o.settings = s
	})
}

// WithFailureCodes sets the failure codes for the circuit breaker.
func WithFailureCodes(cs ...codes.Code) Option {
	return optionFunc(func(o *opts) {
		o.failureCodes = make(map[codes.Code]struct{}, len(cs))
		for _, c := range cs {
			o.failureCodes[c] = struct{}{}
		}
	})
}

type registry struct {
	opts     *opts
	breakers sync.Map
}

// UnaryClientInterceptor returns a gRPC unary client interceptor that uses a circuit breaker to protect the client.
func UnaryClientInterceptor(options ...Option) grpc.UnaryClientInterceptor {
	o := defaultOpts()
	for _, option := range options {
		option.apply(o)
	}

	r := &registry{opts: o}

	return func(ctx context.Context, fullMethod string, req, resp any, conn *grpc.ClientConn, invoker grpc.UnaryInvoker, callOpts ...grpc.CallOption) error {
		cb := r.get(fullMethod)
		_, err := cb.Execute(func() (any, error) {
			return nil, invoker(ctx, fullMethod, req, resp, conn, callOpts...)
		})
		if err != nil {
			if errors.Is(err, breaker.ErrOpenState) || errors.Is(err, breaker.ErrTooManyRequests) {
				return status.Error(codes.Unavailable, err.Error())
			}

			return err
		}
		return nil
	}
}

func (r *registry) get(fullMethod string) *breaker.CircuitBreaker {
	if cb, ok := r.breakers.Load(fullMethod); ok {
		return cb.(*breaker.CircuitBreaker)
	}

	s := r.opts.settings
	s.Name = fullMethod

	failureCodes := r.opts.failureCodes
	s.IsSuccessful = func(err error) bool {
		if err == nil {
			return true
		}

		_, isFailure := failureCodes[status.Code(err)]
		return !isFailure
	}

	cb := breaker.NewCircuitBreaker(s)
	actual, _ := r.breakers.LoadOrStore(fullMethod, cb)
	return actual.(*breaker.CircuitBreaker)
}

func defaultOpts() *opts {
	failureCodes := map[codes.Code]struct{}{
		codes.Unavailable:       {},
		codes.DeadlineExceeded:  {},
		codes.ResourceExhausted: {},
		codes.Internal:          {},
	}

	return &opts{
		failureCodes: failureCodes,
		settings: Settings{
			MaxRequests: 3,
			Interval:    30 * time.Second,
			Timeout:     10 * time.Second,
			ReadyToTrip: func(counts breaker.Counts) bool { return counts.ConsecutiveFailures >= 5 },
		},
	}
}
