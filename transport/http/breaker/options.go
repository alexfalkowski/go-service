package breaker

import (
	"github.com/alexfalkowski/go-service/v2/breaker"
	"github.com/alexfalkowski/go-service/v2/net/http"
)

// Option configures the HTTP circuit breaker RoundTripper.
type Option interface {
	apply(opts *opts)
}

type opts struct {
	settings      Settings
	failureStatus func(code int) bool
}

type optionFunc func(*opts)

func (f optionFunc) apply(o *opts) { f(o) }

// WithSettings configures the circuit breaker settings used for each per-upstream breaker instance.
func WithSettings(s Settings) Option {
	return optionFunc(func(o *opts) { o.settings = s })
}

// WithFailureStatusFunc configures the predicate that classifies an HTTP response status code as a failure
// for breaker accounting.
func WithFailureStatusFunc(f func(code int) bool) Option {
	return optionFunc(func(o *opts) { o.failureStatus = f })
}

// WithFailureStatuses configures a fixed set of HTTP status codes that are treated as failures
// for breaker accounting.
func WithFailureStatuses(statuses ...int) Option {
	set := make(map[int]struct{}, len(statuses))
	for _, s := range statuses {
		set[s] = struct{}{}
	}

	return WithFailureStatusFunc(func(code int) bool {
		_, ok := set[code]
		return ok
	})
}

func defaultOpts() *opts {
	return &opts{
		failureStatus: func(code int) bool {
			return code >= http.StatusInternalServerError || code == http.StatusTooManyRequests
		},
		settings: breaker.DefaultSettings,
	}
}
