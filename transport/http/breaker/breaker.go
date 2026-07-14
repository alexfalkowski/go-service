package breaker

import (
	"cmp"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/transport/breaker"
	"github.com/alexfalkowski/go-sync"
)

// Counts is an alias for [github.com/alexfalkowski/go-service/v2/transport/breaker.Counts].
//
// Counts is used by [github.com/sony/gobreaker.Settings.ReadyToTrip] to decide whether the breaker should open.
type Counts = breaker.Counts

// ErrOpenState is an alias for [github.com/alexfalkowski/go-service/v2/transport/breaker.ErrOpenState].
//
// It is returned by [github.com/sony/gobreaker.CircuitBreaker.Execute] when the breaker is open.
var ErrOpenState = breaker.ErrOpenState

// ErrTooManyRequests is an alias for [github.com/alexfalkowski/go-service/v2/transport/breaker.ErrTooManyRequests].
//
// It is returned by [github.com/sony/gobreaker.CircuitBreaker.Execute] when the breaker is half-open and MaxRequests
// would be exceeded.
var ErrTooManyRequests = breaker.ErrTooManyRequests

// Settings is an alias for [github.com/alexfalkowski/go-service/v2/transport/breaker.Settings].
//
// It is re-exported from this package so callers can configure circuit breaker behavior (trip thresholds,
// timeouts, half-open probing, etc.) without importing the lower-level breaker package directly.
type Settings = breaker.Settings

// NewRoundTripper constructs an HTTP RoundTripper guarded by circuit breakers.
//
// The returned *[RoundTripper] wraps the provided base transport (`hrt`) and executes each request through a
// circuit breaker keyed by the request destination.
//
// # Breaker scope
//
// A separate circuit breaker is maintained per request key: "<METHOD> <HOST>". The host is derived from
// `req.URL.Host`, falling back to `req.Host` (and finally "unknown"). This isolates failures per upstream.
// Breaker keys are retained for the lifetime of the RoundTripper, so this wrapper is intended for a small,
// bounded set of service-to-service upstream hosts rather than arbitrary user-supplied destinations.
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
// Defaults: HTTP status codes >= 500 or 429 are treated as failures (see `defaultOpts` and [WithFailureStatusFunc]).
func NewRoundTripper(hrt http.RoundTripper, options ...Option) *RoundTripper {
	o := defaultOpts()
	for _, option := range options {
		option.apply(o)
	}

	return &RoundTripper{opts: o, RoundTripper: hrt, breakers: sync.NewMap[string, *breaker.CircuitBreaker]()}
}

// RoundTripper wraps an underlying [http.RoundTripper] and applies circuit breaking.
//
// Breakers are cached per request key (method + host) so each upstream is isolated. Breakers are created lazily
// on first use and then reused for subsequent requests to the same key.
// The cache does not evict entries; callers should use this wrapper with bounded upstream host sets.
//
// Use [NewRoundTripper] to construct instances with the desired settings and failure classification behavior.
type RoundTripper struct {
	http.RoundTripper
	opts     *opts
	breakers *sync.Map[string, *breaker.CircuitBreaker]
}

// RoundTrip executes the request guarded by a circuit breaker.
//
// If the breaker is open (or half-open and MaxRequests would be exceeded), the underlying breaker may reject
// the call and return an error (for example [breaker.ErrOpenState] or [breaker.ErrTooManyRequests]).
//
// Failure accounting:
//   - Transport errors are counted as failures.
//   - Responses that match the configured failure-status predicate are counted as failures for breaker
//     accounting, but the response is still returned to the caller with a nil error.
func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return http.ClosingRoundTripper(r.roundTrip).RoundTrip(req)
}

func (r *RoundTripper) roundTrip(req *http.Request) (*http.Response, error, bool) {
	cb := r.get(req)
	result, err := cb.Execute(func() (any, error) {
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
		if re, ok := errors.AsType[responseError](err); ok {
			return re.resp, nil, false
		}

		return nil, localRejectionError(err), errors.Is(err, breaker.ErrOpenState) || errors.Is(err, breaker.ErrTooManyRequests)
	}

	return result.(*http.Response), nil, false
}

func localRejectionError(err error) error {
	if errors.Is(err, breaker.ErrTooManyRequests) {
		return status.LocalError(status.SafeError(http.StatusTooManyRequests, err))
	}
	if errors.Is(err, breaker.ErrOpenState) {
		return status.LocalError(status.ServiceUnavailableError(err))
	}

	return err
}

func (r *RoundTripper) get(req *http.Request) *breaker.CircuitBreaker {
	key := requestKey(req)
	if cb, ok := r.breakers.Load(key); ok {
		return cb
	}

	settings := r.opts.settings
	settings.Name = key
	isSuccessful := settings.IsSuccessful
	settings.IsSuccessful = func(err error) bool {
		if err == nil {
			return true
		}
		if _, ok := errors.AsType[responseError](err); ok {
			return false
		}
		if errors.Is(err, context.Canceled) {
			return true
		}
		if isSuccessful != nil {
			return isSuccessful(err)
		}

		return false
	}

	cb := breaker.NewCircuitBreaker(settings)
	actual, _ := r.breakers.LoadOrStore(key, cb)
	return actual
}

func requestKey(req *http.Request) string {
	// Breaker keys are intentionally stable per method and upstream host. They
	// are retained for the RoundTripper lifetime, so callers should avoid using
	// breaker-enabled clients for arbitrary high-cardinality destinations.
	return strings.Join(" ", req.Method, cmp.Or(req.URL.Host, req.Host, "unknown"))
}
