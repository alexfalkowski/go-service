package token

import (
	"net/http"

	"github.com/alexfalkowski/go-service/errors"
	"github.com/alexfalkowski/go-service/net/http/status"
	"github.com/alexfalkowski/go-service/token"
	"github.com/alexfalkowski/go-service/transport/header"
	"github.com/alexfalkowski/go-service/transport/strings"
	tt "github.com/alexfalkowski/go-service/transport/token"
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

	ctx, err := tt.Verify(ctx, h.verifier)
	if err != nil {
		err := status.FromError(http.StatusUnauthorized, errors.Prefix("token", err))
		status.WriteError(res, err)

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
	ctx, token, err := r.gen.Generate(req.Context())
	if err != nil {
		return nil, err
	}

	if len(token) == 0 {
		return nil, header.ErrInvalidAuthorization
	}

	req = req.WithContext(ctx)
	req.Header.Add("Authorization", header.BearerAuthorization+" "+string(token))

	return r.RoundTripper.RoundTrip(req)
}
