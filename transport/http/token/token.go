package token

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/crypto/ed25519"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/token"
	"github.com/alexfalkowski/go-service/v2/token/access"
	"github.com/alexfalkowski/go-service/v2/transport/header"
	"github.com/alexfalkowski/go-service/v2/transport/strings"
)

// NewAccessController for HTTP.
func NewAccessController(cfg *token.Config) (AccessController, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}
	return access.NewController(cfg.Access)
}

// AccessController is an alias for access.Controller.
type AccessController access.Controller

// NewToken for HTTP.
func NewToken(name env.Name, cfg *token.Config, fs *os.FS, sig *ed25519.Signer, ver *ed25519.Verifier, gen id.Generator) *Token {
	if !cfg.IsEnabled() {
		return nil
	}
	return &Token{Token: token.NewToken(name, cfg, fs, sig, ver, gen)}
}

// Token for HTTP.
type Token struct {
	*token.Token
}

// NewVerifier for HTTP.
func NewVerifier(token *Token) Verifier {
	if token != nil {
		return token
	}
	return nil
}

// Verifier is an alias token.Verifier.
type Verifier token.Verifier

// NewHandler for token.
func NewHandler(id env.UserID, verifier Verifier) *Handler {
	return &Handler{id: id, verifier: verifier}
}

// Handler for token.
type Handler struct {
	verifier Verifier
	id       env.UserID
}

// ServeHTTP verifies the request token and stores the subject in the context.
func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if strings.IsIgnorable(req.URL.Path) {
		next(res, req)
		return
	}

	ctx := req.Context()
	auth := meta.Authorization(ctx).Value()

	sub, err := h.verifier.Verify(strings.Bytes(auth), req.URL.Path)
	if err != nil {
		status.WriteError(ctx, res, status.UnauthorizedError(err))
		return
	}

	ctx = meta.WithUserID(ctx, meta.Ignored(sub))
	next(res, req.WithContext(ctx))
}

// NewGenerator for token.
func NewGenerator(token *Token) Generator {
	if token != nil {
		return token
	}
	return nil
}

// Generator is an alias token.Generator.
type Generator token.Generator

// NewRoundTripper for token.
func NewRoundTripper(id env.UserID, generator Generator, hrt http.RoundTripper) *RoundTripper {
	return &RoundTripper{RoundTripper: hrt, id: id, generator: generator}
}

// RoundTripper for token.
type RoundTripper struct {
	http.RoundTripper
	generator Generator
	id        env.UserID
}

// RoundTrip adds an Authorization header using a generated token.
func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	token, err := r.generator.Generate(req.URL.Path, r.id.String())
	if err != nil {
		return nil, status.UnauthorizedError(err)
	}
	if len(token) == 0 {
		return nil, status.UnauthorizedError(header.ErrInvalidAuthorization)
	}

	req.Header.Add(
		"Authorization",
		strings.Join(strings.Space, header.BearerAuthorization, bytes.String(token)),
	)
	return r.RoundTripper.RoundTrip(req)
}
