package token

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/net/header"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/token"
	"github.com/alexfalkowski/go-service/v2/token/access"
)

// NewToken constructs a token service when token auth is enabled.
//
// The returned service is responsible for generating and verifying tokens according to cfg (for example,
// JWT/PASETO/SSH token kinds as configured by the underlying token package).
//
// If cfg is disabled, it returns nil so callers can treat token auth as not configured.
func NewToken(cfg *token.Config, fs *os.FS, gen id.Generator) *Token {
	if !cfg.IsEnabled() {
		return nil
	}
	return &Token{Token: token.NewToken(cfg, fs, gen)}
}

// Token wraps *[github.com/alexfalkowski/go-service/v2/token.Token] for HTTP transport integration.
//
// It exists so transport-level wiring can keep a distinct type for HTTP token functionality while still
// delegating generation and verification to the underlying token implementation.
type Token struct {
	*token.Token
}

// NewVerifier returns a [Verifier] backed by token.
//
// If token is nil, it returns nil. This pattern allows DI graphs to inject a verifier only when token auth
// is enabled/configured, and to leave verification middleware disabled otherwise.
func NewVerifier(token *Token) Verifier {
	if token != nil {
		return token
	}
	return nil
}

// Verifier is an alias for [token.Verifier].
//
// Verifiers validate Authorization tokens and typically return a "subject" string (the authenticated
// principal) on success.
type Verifier token.Verifier

// NewHandler constructs server-side token verification middleware.
//
// Callers should only install this handler when verifier is non-nil.
func NewHandler(routePolicy *http.RoutePolicy, verifier Verifier) *Handler {
	return &Handler{routePolicy: routePolicy, verifier: verifier}
}

// Handler verifies Authorization headers and injects the verified subject into request metadata.
type Handler struct {
	verifier    Verifier
	routePolicy *http.RoutePolicy
}

// ServeHTTP verifies the request Authorization token and stores the verified subject in the context.
//
// Registered operation paths (health/metrics/etc.) and registered unauthenticated routes bypass verification.
//
// The handler expects an Authorization value to be available in the request context (typically injected by
// [github.com/alexfalkowski/go-service/v2/net/http/meta.Handler]). It verifies the token using verifier, scoping verification to the request method and path.
//
// Behavior:
//   - If verification fails, it writes an HTTP 401 error response and does not call next.
//   - If verification succeeds, it stores the verified subject as the user id in the request context and calls next.
//
// Callers should only install this handler when verifier is non-nil.
func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if h.routePolicy.IsOperation(req) || h.routePolicy.IsUnauthenticated(req) {
		next(res, req)
		return
	}

	ctx := req.Context()
	auth := meta.Authorization(ctx).Value()

	sub, err := h.verifier.Verify(strings.Bytes(auth), audience(req))
	if err != nil {
		_ = status.WriteError(req.Context(), res, status.UnauthorizedError(err))
		return
	}

	ctx = meta.WithAttributes(ctx, meta.WithUserID(meta.Ignored(sub)))
	next(res, req.WithContext(ctx))
}

// NewAccessHandler constructs server-side access-control middleware.
//
// Callers should only install this handler when controller is non-nil.
func NewAccessHandler(routePolicy *http.RoutePolicy, controller access.Controller) *AccessHandler {
	return &AccessHandler{routePolicy: routePolicy, controller: controller}
}

// AccessHandler enforces access policy for token-authenticated requests.
type AccessHandler struct {
	controller  access.Controller
	routePolicy *http.RoutePolicy
}

// ServeHTTP checks the verified user id against the configured access policy.
//
// Registered operation paths (health/metrics/etc.) and registered unauthenticated routes bypass access
// control. For application paths, a missing verified user id is treated as unauthenticated, a policy denial
// is written as HTTP 403, and policy evaluation errors are written as HTTP 500.
func (h *AccessHandler) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if h.routePolicy.IsOperation(req) || h.routePolicy.IsUnauthenticated(req) {
		next(res, req)
		return
	}

	ctx := req.Context()
	if meta.UserID(ctx).IsEmpty() {
		_ = status.WriteError(req.Context(), res, status.UnauthorizedError(header.ErrInvalidAuthorization))
		return
	}

	ok, err := h.controller.HasAccess(ctx)
	if err != nil {
		_ = status.WriteError(req.Context(), res, status.InternalServerError(err))
		return
	}
	if !ok {
		_ = status.WriteError(req.Context(), res, status.SafeError(http.StatusForbidden, access.ErrAccessDenied))
		return
	}

	next(res, req)
}

// NewGenerator returns a [Generator] backed by token.
//
// If token is nil, it returns nil. This pattern allows DI graphs to inject a token generator only when token
// auth is enabled/configured, and to leave client token-injection middleware disabled otherwise.
func NewGenerator(token *Token) Generator {
	if token != nil {
		return token
	}
	return nil
}

// Generator is an alias for [token.Generator].
//
// Generators create Authorization tokens for outbound HTTP requests, typically scoped to the request method
// and path plus a caller identity (user id).
type Generator token.Generator

// NewRoundTripper constructs client-side token injection middleware for HTTP requests.
//
// For each outbound request, the RoundTripper generates a token and sets the request's Authorization
// header to a single `Bearer` value.
//
// Callers should only install this RoundTripper when generator is non-nil.
func NewRoundTripper(id env.UserID, generator Generator, hrt http.RoundTripper) *RoundTripper {
	return &RoundTripper{RoundTripper: hrt, id: id, generator: generator}
}

// RoundTripper wraps an underlying [http.RoundTripper] and adds Authorization headers.
type RoundTripper struct {
	http.RoundTripper
	generator Generator
	id        env.UserID
}

// RoundTrip sets the Authorization header using a generated token.
//
// For each request, it generates a token scoped to the request method and path and the configured user id
// and then writes it as a single `Bearer` token in the Authorization header.
//
// Failure behavior:
//   - If token generation fails, it returns an unauthorized status error.
//   - If token generation returns an empty token, it returns an unauthorized status error with [header.ErrInvalidAuthorization].
//   - If the request is a cross-origin redirect, it returns [http.ErrUseLastResponse] without forwarding
//     credentials to the redirected origin.
func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return http.ClosingRoundTripper(r.roundTrip).RoundTrip(req)
}

func (r *RoundTripper) roundTrip(req *http.Request) (*http.Response, error, bool) {
	if http.IsCrossOriginRedirect(req) {
		return nil, http.ErrUseLastResponse, true
	}

	token, err := r.generator.Generate(audience(req), r.id.String())
	if err != nil {
		return nil, status.UnauthorizedError(err), true
	}
	if len(token) == 0 {
		return nil, status.UnauthorizedError(header.ErrInvalidAuthorization), true
	}

	auth := meta.Ignored(strings.Join(strings.Space, header.BearerAuthorization, bytes.String(token)))
	ctx := meta.WithAttributes(req.Context(), meta.WithAuthorization(auth))

	clonedReq := req.Clone(ctx)
	if clonedReq.Header == nil {
		clonedReq.Header = http.Header{}
	}
	clonedReq.Header.Set("Authorization", auth.Value())

	res, err := r.RoundTripper.RoundTrip(clonedReq)
	return res, err, false
}

func audience(req *http.Request) string {
	return strings.Join(strings.Space, req.Method, req.URL.Path)
}
