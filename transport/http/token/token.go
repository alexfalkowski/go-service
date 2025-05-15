package token

import (
	"net/http"

	"github.com/alexfalkowski/go-service/bytes"
	"github.com/alexfalkowski/go-service/errors"
	"github.com/alexfalkowski/go-service/net/http/status"
	"github.com/alexfalkowski/go-service/strings"
	"github.com/alexfalkowski/go-service/token"
	"github.com/alexfalkowski/go-service/transport/header"
	"github.com/alexfalkowski/go-service/transport/meta"
	ts "github.com/alexfalkowski/go-service/transport/strings"
)

// NewHandler for token.
func NewHandler(verifier token.Verifier) *Handler {
	return &Handler{verifier: verifier}
}

// Handler for token.
type Handler struct {
	verifier token.Verifier
}

func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if ts.IsObservable(req.URL.Path) {
		next(res, req)

		return
	}

	ctx := req.Context()
	token := meta.Authorization(ctx).Value()

	ctx, err := h.verifier.Verify(ctx, []byte(token))
	if err != nil {
		status.WriteError(ctx, res, status.UnauthorizedError(errors.Prefix("token", err)))

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
	req.Header.Add(
		"Authorization",
		strings.Join(" ", header.BearerAuthorization, bytes.String(token)),
	)

	return r.RoundTripper.RoundTrip(req)
}
