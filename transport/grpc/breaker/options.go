package breaker

import (
	"github.com/alexfalkowski/go-service/v2/breaker"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
)

// Option configures the gRPC circuit breaker interceptor.
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
func WithSettings(s Settings) Option {
	return optionFunc(func(o *opts) {
		o.settings = s
	})
}

// WithFailureCodes configures which gRPC status codes are treated as failures for breaker accounting.
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
