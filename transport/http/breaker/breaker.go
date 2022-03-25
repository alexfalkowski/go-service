package breaker

import (
	"net/http"

	breaker "github.com/sony/gobreaker"
)

// NewRoundTripper for breaker.
func NewRoundTripper(hrt http.RoundTripper) *RoundTripper {
	cb := breaker.NewCircuitBreaker(breaker.Settings{})

	return &RoundTripper{cb: cb, RoundTripper: hrt}
}

// RoundTripper for breaker.
type RoundTripper struct {
	cb *breaker.CircuitBreaker
	http.RoundTripper
}

// nolint:forcetypeassert
func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	operation := func() (interface{}, error) {
		return r.RoundTripper.RoundTrip(req)
	}

	resp, err := r.cb.Execute(operation)
	if err != nil {
		return nil, err
	}

	return resp.(*http.Response), nil
}
