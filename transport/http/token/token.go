package token

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/env"
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
func NewHandler(id env.UserID, verifier Verifier) *Handler {
	return &Handler{id: id, verifier: verifier}
}

// Handler for token.
type Handler struct {
	verifier Verifier
	id       env.UserID
}

func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	p := http.Path(req)
	if strings.IsObservable(p) {
		next(res, req)

		return
	}

	ctx := req.Context()
	auth := meta.Authorization(ctx).Value()

	sub, err := h.verifier.Verify(strings.Bytes(auth), p)
	if err != nil {
		status.WriteError(ctx, res, status.UnauthorizedError(err))

		return
	}

	ctx = meta.WithUserID(ctx, meta.Ignored(sub))

	next(res, req.WithContext(ctx))
}

// NewRoundTripper for token.
func NewRoundTripper(id env.UserID, generator token.Generator, hrt http.RoundTripper) *RoundTripper {
	return &RoundTripper{RoundTripper: hrt, id: id, generator: generator}
}

// RoundTripper for token.
type RoundTripper struct {
	http.RoundTripper

	generator Generator
	id        env.UserID
}

func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	p := http.Path(req)

	token, err := r.generator.Generate(p, r.id.String())
	if err != nil {
		return nil, status.UnauthorizedError(err)
	}

	if len(token) == 0 {
		return nil, status.UnauthorizedError(header.ErrInvalidAuthorization)
	}

	req.Header.Add(
		"Authorization",
		strings.Join(" ", header.BearerAuthorization, bytes.String(token)),
	)

	return r.RoundTripper.RoundTrip(req)
}
