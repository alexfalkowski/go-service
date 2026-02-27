package breaker

import (
	"github.com/alexfalkowski/go-service/v2/breaker"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
)

// Option configures the gRPC circuit breaker interceptor returned by `UnaryClientInterceptor`.
//
// Options are applied in the order provided to `UnaryClientInterceptor`. If multiple options configure
// the same field, the last one wins.
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

// WithSettings configures the circuit breaker settings used for each per-method breaker instance.
//
// The settings value is copied into each created breaker, and the interceptor wiring will also set:
//
//   - `Settings.Name` to the `fullMethod`, and
//   - `Settings.IsSuccessful` to treat selected gRPC status codes as failures (see `WithFailureCodes`).
//
// Note: because settings are copied, if your `Settings` contains function fields that close over
// mutable state, ensure that state is safe for concurrent use.
func WithSettings(s Settings) Option {
	return optionFunc(func(o *opts) {
		o.settings = s
	})
}

// WithFailureCodes configures which gRPC status codes are treated as failures for breaker accounting.
//
// If an invocation returns an error whose status code is contained in this set, the breaker counts it as a
// failure. Errors with status codes not in this set are counted as successes.
func WithFailureCodes(cs ...codes.Code) Option {
	return optionFunc(func(o *opts) {
		o.failureCodes = make(map[codes.Code]struct{}, len(cs))
		for _, c := range cs {
			o.failureCodes[c] = struct{}{}
		}
	})
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
		settings:     breaker.DefaultSettings,
	}
}
