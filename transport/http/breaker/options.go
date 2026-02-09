package breaker

import (
	"github.com/alexfalkowski/go-service/v2/breaker"
	"github.com/alexfalkowski/go-service/v2/net/http"
)

// Option interface for configuring the circuit breaker.
type Option interface {
	apply(opts *opts)
}

type opts struct {
	settings      Settings
	failureStatus func(code int) bool
}

type optionFunc func(*opts)

func (f optionFunc) apply(o *opts) { f(o) }

// WithSettings for configuring the circuit breaker.
func WithSettings(s Settings) Option {
	return optionFunc(func(o *opts) { o.settings = s })
}

// WithFailureStatusFunc for configuring the circuit breaker.
func WithFailureStatusFunc(f func(code int) bool) Option {
	return optionFunc(func(o *opts) { o.failureStatus = f })
}

// WithFailureStatuses for configuring the circuit breaker.
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
