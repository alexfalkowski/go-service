package token

import (
	"fmt"
	"net/http"

	"github.com/alexfalkowski/go-service/errors"
	nh "github.com/alexfalkowski/go-service/net/http"
	"github.com/alexfalkowski/go-service/security/header"
	"github.com/alexfalkowski/go-service/security/token"
	st "github.com/alexfalkowski/go-service/transport/security/token"
	"github.com/alexfalkowski/go-service/transport/strings"
)

// Handler for token.
type Handler struct {
	verifier token.Verifier
}

// NewHandler for token.
func NewHandler(verifier token.Verifier) *Handler {
	return &Handler{verifier: verifier}
}

func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if strings.IsObservable(req.URL.Path) {
		next(res, req)

		return
	}

	ctx := req.Context()

	ctx, err := st.Verify(ctx, h.verifier)
	if err != nil {
		nh.WriteError(req.Context(), res, errors.Prefix("verify token", err), http.StatusUnauthorized)

		return
	}

	next(res, req.WithContext(ctx))
}

// NewRoundTripper for token.
func NewRoundTripper(gen token.Generator, hrt http.RoundTripper) *RoundTripper {
	return &RoundTripper{gen: gen, RoundTripper: hrt}
}

// RoundTripper for token.
type RoundTripper struct {
	gen token.Generator

	http.RoundTripper
}

func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx, t, err := r.gen.Generate(req.Context())
	if err != nil {
		return nil, err
	}

	if len(t) == 0 {
		return nil, header.ErrInvalidAuthorization
	}

	req = req.WithContext(ctx)
	req.Header.Add("Authorization", fmt.Sprintf("%s %s", header.BearerAuthorization, string(t)))

	return r.RoundTripper.RoundTrip(req)
}
