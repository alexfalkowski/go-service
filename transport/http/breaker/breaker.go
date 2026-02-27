package breaker

import (
	"cmp"

	"github.com/alexfalkowski/go-service/v2/breaker"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/sync"
	"github.com/alexfalkowski/go-service/v2/transport/strings"
)

// Settings is an alias for `github.com/alexfalkowski/go-service/v2/breaker.Settings`.
//
// It is re-exported from this package so callers can configure circuit breaker behavior (trip thresholds,
// timeouts, half-open probing, etc.) without importing the lower-level breaker package directly.
type Settings = breaker.Settings

// NewRoundTripper constructs an HTTP RoundTripper guarded by circuit breakers.
//
// The returned `*RoundTripper` wraps the provided base transport (`hrt`) and executes each request through a
// circuit breaker keyed by the request destination.
//
// # Breaker scope
//
// A separate circuit breaker is maintained per request key: "<METHOD> <HOST>". The host is derived from
// `req.URL.Host`, falling back to `req.Host` (and finally "unknown"). This isolates failures per upstream.
//
// # Failure accounting and error semantics
//
// The breaker executes the underlying transport call and classifies the outcome for breaker accounting:
//
//   - Transport errors (i.e., the underlying RoundTripper returns a non-nil error) are counted as failures.
//   - HTTP responses whose `StatusCode` matches the configured failure-status predicate are also counted as failures.
//
// Important: when a response status code is treated as a failure for breaker accounting, this RoundTripper still
// returns the response to the caller with a nil error. This means:
//
//   - Your application logic continues to be driven by the HTTP response status/body, and
//   - The breaker still "learns" that the upstream is unhealthy and may open accordingly.
//
// Defaults: HTTP status codes >= 500 or 429 are treated as failures (see `defaultOpts` and `WithFailureStatusFunc`).
func NewRoundTripper(hrt http.RoundTripper, options ...Option) *RoundTripper {
	o := defaultOpts()
	for _, option := range options {
		option.apply(o)
	}

	return &RoundTripper{opts: o, RoundTripper: hrt, breakers: sync.NewMap[string, *breaker.CircuitBreaker]()}
}

// RoundTripper wraps an underlying `http.RoundTripper` and applies circuit breaking.
//
// Breakers are cached per request key (method + host) so each upstream is isolated. Breakers are created lazily
// on first use and then reused for subsequent requests to the same key.
//
// Use `NewRoundTripper` to construct instances with the desired settings and failure classification behavior.
type RoundTripper struct {
	http.RoundTripper
	opts     *opts
	breakers sync.Map[string, *breaker.CircuitBreaker]
}

// RoundTrip executes the request guarded by a circuit breaker.
//
// If the breaker is open (or half-open and MaxRequests would be exceeded), the underlying breaker may reject
// the call and return an error (for example `breaker.ErrOpenState` or `breaker.ErrTooManyRequests`).
//
// Failure accounting:
//   - Transport errors are counted as failures.
//   - Responses that match the configured failure-status predicate are counted as failures for breaker
//     accounting, but the response is still returned to the caller with a nil error.
func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	cb := r.get(req)
	v, err := cb.Execute(func() (any, error) {
		resp, err := r.RoundTripper.RoundTrip(req)
		if err != nil {
			return nil, err
		}

		if r.opts.failureStatus(resp.StatusCode) {
			return nil, responseError{resp: resp}
		}

		return resp, nil
	})
	if err != nil {
		var re responseError
		if errors.As(err, &re) {
			return re.resp, nil
		}

		return nil, err
	}
	return v.(*http.Response), nil
}

func (r *RoundTripper) get(req *http.Request) *breaker.CircuitBreaker {
	key := requestKey(req)
	if cb, ok := r.breakers.Load(key); ok {
		return cb
	}

	s := r.opts.settings
	s.Name = key
	s.IsSuccessful = func(err error) bool {
		if err != nil {
			var re responseError
			return !errors.As(err, &re)
		}

		return true
	}

	cb := breaker.NewCircuitBreaker(s)
	actual, _ := r.breakers.LoadOrStore(key, cb)
	return actual
}

func requestKey(req *http.Request) string {
	return strings.Join(" ", req.Method, cmp.Or(req.URL.Host, req.Host, "unknown"))
}
