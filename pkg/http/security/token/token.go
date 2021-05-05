package token

import (
	"fmt"
	"net/http"

	"github.com/alexfalkowski/go-service/pkg/security/token"
)

// NewRoundTripper for token.
func NewRoundTripper(gen token.Generator, hrt http.RoundTripper) *RoundTripper {
	return &RoundTripper{gen: gen, RoundTripper: hrt}
}

// RoundTripper for zap.
type RoundTripper struct {
	gen token.Generator

	http.RoundTripper
}

func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	token, err := r.gen.Generate()
	if err != nil {
		return nil, err
	}

	if len(token) == 0 {
		return r.RoundTripper.RoundTrip(req)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", string(token)))

	return r.RoundTripper.RoundTrip(req)
}
