package breaker

import (
	"github.com/alexfalkowski/go-service/v2/breaker"
	"github.com/alexfalkowski/go-service/v2/net/http"
)

// Option configures the HTTP circuit breaker `RoundTripper` created by `NewRoundTripper`.
//
// Options are applied in the order provided to `NewRoundTripper`. If multiple options configure
// the same field, the last one wins.
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
//
// The settings value is copied into each created breaker, and `NewRoundTripper` will also set:
//
//   - `Settings.Name` to the request key ("<METHOD> <HOST>"), and
//   - `Settings.IsSuccessful` to ensure responses classified as failures by this package's
//     failure-status predicate are counted as failures for breaker accounting.
//
// Note: because settings are copied, if your `Settings` contains function fields that close over
// mutable state, ensure that state is safe for concurrent use.
func WithSettings(s Settings) Option {
	return optionFunc(func(o *opts) { o.settings = s })
}

// WithFailureStatusFunc configures the predicate that classifies an HTTP response status code as a failure
// for breaker accounting.
//
// When the predicate returns true for a response status code, the breaker counts the call as a failure,
// but the `RoundTripper` still returns the original `*http.Response` to the caller with a nil error.
// This decouples breaker health tracking from application-level HTTP response handling.
func WithFailureStatusFunc(f func(code int) bool) Option {
	return optionFunc(func(o *opts) { o.failureStatus = f })
}

// WithFailureStatuses configures a fixed set of HTTP status codes that are treated as failures
// for breaker accounting.
//
// This is a convenience wrapper over `WithFailureStatusFunc`. The provided codes are stored in a set and
// membership is checked for each response.
//
// Example: treat 502/503/504 as failures:
//
//	WithFailureStatuses(502, 503, 504)
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
