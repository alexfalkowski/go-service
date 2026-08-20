package breaker

import (
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/transport/breaker"
)

// Option configures the gRPC circuit breaker client returned by [NewClient].
//
// Options are applied in the order provided to [NewClient]. If multiple options configure
// the same field, the last one wins.
type Option interface {
	apply(options *options)
}

type options struct {
	settings     Settings
	failureCodes map[codes.Code]struct{}
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) {
	f(o)
}

// WithSettings configures the circuit breaker settings used for each per-method breaker instance.
//
// The settings value is copied into each created breaker, and the interceptor wiring will also set:
//
//   - [github.com/sony/gobreaker.Settings.Name] to the `fullMethod`, and
//   - [github.com/sony/gobreaker.Settings.IsSuccessful] to treat selected gRPC status codes as failures (see [WithFailureCodes]).
//
// Note: because settings are copied, if your [Settings] contains function fields that close over
// mutable state, ensure that state is safe for concurrent use.
func WithSettings(s Settings) Option {
	return optionFunc(func(o *options) {
		o.settings = s
	})
}

// WithFailureCodes configures which gRPC status codes are treated as failures for breaker accounting.
//
// If an invocation returns an error whose status code is contained in this set, the breaker counts it as a
// failure. Errors with status codes not in this set are counted as successes. [codes.Canceled] is always
// counted as a success regardless of this set, matching the caller-cancellation handling in the HTTP client
// breaker: a caller aborting an in-flight call should not trip the breaker against a healthy upstream. An
// already-expired caller deadline bypasses breaker accounting before invocation for the same reason.
func WithFailureCodes(cs ...codes.Code) Option {
	return optionFunc(func(o *options) {
		o.failureCodes = make(map[codes.Code]struct{}, len(cs))
		for _, c := range cs {
			o.failureCodes[c] = struct{}{}
		}
	})
}

func defaultOptions() *options {
	failureCodes := map[codes.Code]struct{}{
		codes.Unavailable:       {},
		codes.DeadlineExceeded:  {},
		codes.ResourceExhausted: {},
		codes.Internal:          {},
	}

	return &options{
		failureCodes: failureCodes,
		settings:     breaker.DefaultSettings,
	}
}
