package breaker

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/sync"
	"github.com/alexfalkowski/go-service/v2/time"
	breaker "github.com/sony/gobreaker"
)

// ErrInvalidResponse is returned when the response is invalid.
var ErrInvalidResponse = errors.New("breaker: invalid response")

// Option interface for configuring the circuit breaker.
type Option interface {
	apply(opts *opts)
}

type opts struct {
	settings      breaker.Settings
	failureStatus func(code int) bool
}

type optionFunc func(*opts)

func (f optionFunc) apply(o *opts) { f(o) }

// WithSettings for configuring the circuit breaker.
func WithSettings(s breaker.Settings) Option {
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

// NewRoundTripper for breaker.
func NewRoundTripper(hrt http.RoundTripper, options ...Option) *RoundTripper {
	o := defaultOpts()
	for _, option := range options {
		option.apply(o)
	}

	return &RoundTripper{opts: o, RoundTripper: hrt}
}

// RoundTripper for breaker.
type RoundTripper struct {
	http.RoundTripper
	opts     *opts
	breakers sync.Map
}

func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	cb := r.get(req)

	v, err := cb.Execute(func() (any, error) {
		resp, err := r.RoundTripper.RoundTrip(req)
		if err != nil {
			return nil, err
		}

		if r.opts.failureStatus(resp.StatusCode) {
			return resp, responseError{resp: resp}
		}

		return resp, nil
	})
	if err == nil {
		resp, ok := v.(*http.Response)
		if !ok {
			return nil, ErrInvalidResponse
		}
		return resp, nil
	}

	var re responseError
	if errors.As(err, &re) {
		return re.resp, nil
	}

	return nil, err
}

func (r *RoundTripper) get(req *http.Request) *breaker.CircuitBreaker {
	key := requestKey(req)
	if cb, ok := r.breakers.Load(key); ok {
		return cb.(*breaker.CircuitBreaker)
	}

	s := r.opts.settings
	s.Name = key
	s.IsSuccessful = func(err error) bool {
		if err == nil {
			return true
		}

		if errors.Is(err, context.Canceled) {
			return true
		}

		var re responseError
		return !errors.As(err, &re)
	}

	cb := breaker.NewCircuitBreaker(s)
	actual, _ := r.breakers.LoadOrStore(key, cb)
	return actual.(*breaker.CircuitBreaker)
}

func requestKey(req *http.Request) string {
	host := req.URL.Host
	if host == "" {
		host = req.Host
	}
	if host == "" {
		host = "unknown"
	}

	return req.Method + " " + host
}

func defaultOpts() *opts {
	return &opts{
		failureStatus: func(code int) bool {
			return code >= http.StatusInternalServerError || code == http.StatusTooManyRequests
		},
		settings: breaker.Settings{
			MaxRequests: 3,
			Interval:    30 * time.Second,
			Timeout:     10 * time.Second,
			ReadyToTrip: func(counts breaker.Counts) bool { return counts.ConsecutiveFailures >= 5 },
		},
	}
}

type responseError struct {
	resp *http.Response
}

func (e responseError) Error() string {
	if e.resp == nil {
		return "breaker: failure response"
	}
	return "breaker: failure response (status " + e.resp.Status + ")"
}
