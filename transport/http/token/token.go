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

// NewAccessController constructs an access controller when token auth is enabled.
//
// The controller is responsible for evaluating access rules associated with token-authenticated subjects.
//
// If cfg is disabled, it returns (nil, nil) so callers can treat access control as not configured.
func NewAccessController(cfg *token.Config) (AccessController, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}
	return access.NewController(cfg.Access)
}

// AccessController is an alias for `token/access.Controller`.
//
// It is exposed from this package so callers can refer to the access controller type from the HTTP token
// integration layer.
type AccessController access.Controller

// NewToken constructs a token service when token auth is enabled.
//
// The returned service is responsible for generating and verifying tokens according to cfg (for example,
// JWT/PASETO/SSH token kinds as configured by the underlying token package).
//
// If cfg is disabled, it returns nil so callers can treat token auth as not configured.
func NewToken(name env.Name, cfg *token.Config, fs *os.FS, sig *ed25519.Signer, ver *ed25519.Verifier, gen id.Generator) *Token {
	if !cfg.IsEnabled() {
		return nil
	}
	return &Token{Token: token.NewToken(name, cfg, fs, sig, ver, gen)}
}

// Token wraps `*token.Token` for HTTP transport integration.
//
// It exists so transport-level wiring can keep a distinct type for HTTP token functionality while still
// delegating generation and verification to the underlying token implementation.
type Token struct {
	*token.Token
}

// NewVerifier returns a `Verifier` backed by token.
//
// If token is nil, it returns nil. This pattern allows DI graphs to inject a verifier only when token auth
// is enabled/configured, and to leave verification middleware disabled otherwise.
func NewVerifier(token *Token) Verifier {
	if token != nil {
		return token
	}
	return nil
}

// Verifier is an alias for `token.Verifier`.
//
// Verifiers validate Authorization tokens and typically return a "subject" string (the authenticated
// principal) on success.
type Verifier token.Verifier

// NewHandler constructs server-side token verification middleware.
//
// Callers should only install this handler when verifier is non-nil.
func NewHandler(id env.UserID, verifier Verifier) *Handler {
	return &Handler{id: id, verifier: verifier}
}

// Handler verifies Authorization headers and injects the verified subject into request metadata.
type Handler struct {
	verifier Verifier
	id       env.UserID
}

// ServeHTTP verifies the request Authorization token and stores the verified subject in the context.
//
// Ignorable paths (health/metrics/etc.) bypass verification (see `transport/strings.IsIgnorable`).
//
// The handler expects an Authorization value to be available in the request context (typically injected by
// `transport/http/meta.Handler`). It verifies the token using verifier, scoping verification to the request path.
//
// Behavior:
//   - If verification fails, it writes an HTTP 401 error response and does not call next.
//   - If verification succeeds, it stores the verified subject as the user id in the request context and calls next.
//
// Callers should only install this handler when verifier is non-nil.
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

// NewGenerator returns a `Generator` backed by token.
//
// If token is nil, it returns nil. This pattern allows DI graphs to inject a token generator only when token
// auth is enabled/configured, and to leave client token-injection middleware disabled otherwise.
func NewGenerator(token *Token) Generator {
	if token != nil {
		return token
	}
	return nil
}

// Generator is an alias for `token.Generator`.
//
// Generators create Authorization tokens for outbound HTTP requests, typically scoped to the request path
// and a caller identity (user id).
type Generator token.Generator

// NewRoundTripper constructs client-side token injection middleware for HTTP requests.
//
// For each outbound request, the RoundTripper generates a token and adds it to the request's Authorization
// header using the `Bearer` scheme.
//
// Callers should only install this RoundTripper when generator is non-nil.
func NewRoundTripper(id env.UserID, generator Generator, hrt http.RoundTripper) *RoundTripper {
	return &RoundTripper{RoundTripper: hrt, id: id, generator: generator}
}

// RoundTripper wraps an underlying `http.RoundTripper` and adds Authorization headers.
type RoundTripper struct {
	http.RoundTripper
	generator Generator
	id        env.UserID
}

// RoundTrip adds an Authorization header using a generated token.
//
// For each request, it generates a token scoped to the request path and the configured user id and then
// adds it as a `Bearer` token in the Authorization header.
//
// Failure behavior:
//   - If token generation fails, it returns an unauthorized status error.
//   - If token generation returns an empty token, it returns an unauthorized status error with `header.ErrInvalidAuthorization`.
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
