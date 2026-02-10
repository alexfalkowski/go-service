package breaker

import (
	"cmp"

	"github.com/alexfalkowski/go-service/v2/breaker"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/sync"
	"github.com/alexfalkowski/go-service/v2/transport/strings"
)

// Settings is an alias for breaker.Settings.
type Settings = breaker.Settings

// NewRoundTripper returns an HTTP RoundTripper guarded by circuit breakers.
//
// A separate circuit breaker is maintained per request key (method + host).
// By default, HTTP responses with status codes >= 500 or 429 are treated as failures.
func NewRoundTripper(hrt http.RoundTripper, options ...Option) *RoundTripper {
	o := defaultOpts()
	for _, option := range options {
		option.apply(o)
	}

	return &RoundTripper{opts: o, RoundTripper: hrt, breakers: sync.NewMap[string, *breaker.CircuitBreaker]()}
}

// RoundTripper wraps an underlying http.RoundTripper and applies circuit breaking.
//
// Circuit breakers are cached per request key (method + host) so each upstream is isolated.
type RoundTripper struct {
	http.RoundTripper
	opts     *opts
	breakers sync.Map[string, *breaker.CircuitBreaker]
}

// RoundTrip executes the request guarded by a circuit breaker.
//
// Transport errors are counted as failures.
// Responses that match the configured failure-status predicate are treated as failures for breaker accounting,
// but the response is returned to the caller (with a nil error).
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
