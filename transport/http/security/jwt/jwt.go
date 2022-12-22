package jwt

import (
	"fmt"
	"net/http"

	"github.com/alexfalkowski/go-service/security/header"
	"github.com/alexfalkowski/go-service/security/jwt"
)

// NewRoundTripper for token.
func NewRoundTripper(gen jwt.Generator, hrt http.RoundTripper) *RoundTripper {
	return &RoundTripper{gen: gen, RoundTripper: hrt}
}

// RoundTripper for zap.
type RoundTripper struct {
	gen jwt.Generator

	http.RoundTripper
}

func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	t, err := r.gen.Generate(req.Context())
	if err != nil {
		return nil, err
	}

	if len(t) == 0 {
		return nil, header.ErrInvalidAuthorization
	}

	req.Header.Add("Authorization", fmt.Sprintf("%s %s", header.BearerAuthorization, string(t)))

	return r.RoundTripper.RoundTrip(req)
}
