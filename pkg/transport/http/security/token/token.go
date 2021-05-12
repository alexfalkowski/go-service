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
	t, err := r.gen.Generate(req.Context())
	if err != nil {
		return nil, err
	}

	if len(t) == 0 {
		return nil, token.ErrMissingToken
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", string(t)))

	return r.RoundTripper.RoundTrip(req)
}
