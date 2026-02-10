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

// NewAccessController returns an access controller when token auth is enabled.
//
// If cfg is disabled, it returns (nil, nil).
func NewAccessController(cfg *token.Config) (AccessController, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}
	return access.NewController(cfg.Access)
}

// AccessController is an alias for access.Controller.
type AccessController access.Controller

// NewToken returns a token service when token auth is enabled.
//
// If cfg is disabled, it returns nil.
func NewToken(name env.Name, cfg *token.Config, fs *os.FS, sig *ed25519.Signer, ver *ed25519.Verifier, gen id.Generator) *Token {
	if !cfg.IsEnabled() {
		return nil
	}
	return &Token{Token: token.NewToken(name, cfg, fs, sig, ver, gen)}
}

// Token wraps token.Token for HTTP transport integration.
type Token struct {
	*token.Token
}

// NewVerifier returns a Verifier backed by token.
//
// If token is nil, it returns nil.
func NewVerifier(token *Token) Verifier {
	if token != nil {
		return token
	}
	return nil
}

// Verifier is an alias for token.Verifier.
type Verifier token.Verifier

// NewHandler constructs a token verification handler.
func NewHandler(id env.UserID, verifier Verifier) *Handler {
	return &Handler{id: id, verifier: verifier}
}

// Handler for token.
type Handler struct {
	verifier Verifier
	id       env.UserID
}

// ServeHTTP verifies the request Authorization token and stores the subject in the context.
//
// Requests with ignorable paths bypass verification.
// On verification failure, it writes an unauthorized error response.
// On success, it stores the verified subject as the user id in the request context and calls next.
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

// NewGenerator returns a Generator backed by token.
//
// If token is nil, it returns nil.
func NewGenerator(token *Token) Generator {
	if token != nil {
		return token
	}
	return nil
}

// Generator is an alias for token.Generator.
type Generator token.Generator

// NewRoundTripper constructs an HTTP RoundTripper that adds an Authorization header.
//
// The token is generated per request using generator and the provided user id.
func NewRoundTripper(id env.UserID, generator Generator, hrt http.RoundTripper) *RoundTripper {
	return &RoundTripper{RoundTripper: hrt, id: id, generator: generator}
}

// RoundTripper wraps an underlying http.RoundTripper and adds Authorization headers.
type RoundTripper struct {
	http.RoundTripper
	generator Generator
	id        env.UserID
}

// RoundTrip adds an Authorization header using a generated token.
//
// It generates a token scoped to the request path and the configured user id and then adds it as a
// Bearer token in the Authorization header.
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
