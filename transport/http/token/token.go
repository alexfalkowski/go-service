package token

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/token"
	"github.com/alexfalkowski/go-service/v2/transport/header"
	"github.com/alexfalkowski/go-service/v2/transport/meta"
	"github.com/alexfalkowski/go-service/v2/transport/strings"
)

type (
	// Generator is an alias token.Generator.
	Generator = token.Generator

	// Verifier is an alias token.Verifier.
	Verifier = token.Verifier
)

// NewHandler for token.
func NewHandler(verifier Verifier) *Handler {
	return &Handler{verifier: verifier}
}

// Handler for token.
type Handler struct {
	verifier Verifier
}

func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	p := http.Path(req)
	if strings.IsObservable(p) {
		next(res, req)

		return
	}

	ctx := req.Context()
	auth := meta.Authorization(ctx).Value()

	ctx, err := h.verifier.Verify(ctx, strings.Bytes(auth), token.Options{Path: p})
	if err != nil {
		status.WriteError(ctx, res, status.UnauthorizedError(err))

		return
	}

	next(res, req.WithContext(ctx))
}

// NewRoundTripper for token.
func NewRoundTripper(generator token.Generator, hrt http.RoundTripper) *RoundTripper {
	return &RoundTripper{generator: generator, RoundTripper: hrt}
}

// RoundTripper for token.
type RoundTripper struct {
	generator Generator

	http.RoundTripper
}

func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	p := http.Path(req)

	ctx, token, err := r.generator.Generate(req.Context(), token.Options{Path: p})
	if err != nil {
		return nil, status.UnauthorizedError(err)
	}

	if len(token) == 0 {
		return nil, status.UnauthorizedError(header.ErrInvalidAuthorization)
	}

	req = req.WithContext(ctx)
	req.Header.Add(
		"Authorization",
		strings.Join(" ", header.BearerAuthorization, bytes.String(token)),
	)

	return r.RoundTripper.RoundTrip(req)
}
