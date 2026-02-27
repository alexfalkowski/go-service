package limiter

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/limiter"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/transport/strings"
)

// KeyMap is an alias for `limiter.KeyMap`.
//
// It maps limiter key kinds (for example, "user-agent" or "ip") to functions that derive a rate-limit key
// from the request context.
type KeyMap = limiter.KeyMap

// NewServerLimiter constructs an HTTP server-side rate limiter.
//
// If cfg is disabled, it returns (nil, nil) so callers can treat the limiter as not configured.
//
// The returned limiter is backed by `limiter.NewLimiter` and is registered with the provided lifecycle.
// The `keys` map controls how request contexts are turned into limiter keys (for example, per user-agent).
func NewServerLimiter(lc di.Lifecycle, keys KeyMap, cfg *limiter.Config) (*Server, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}

	limiter, err := limiter.NewLimiter(lc, keys, cfg)
	if err != nil {
		return nil, err
	}

	return &Server{limiter}, nil
}

// Server wraps `*limiter.Limiter` for HTTP server integration.
type Server struct {
	*limiter.Limiter
}

// NewHandler constructs a server-side rate limiting Negroni handler.
//
// Callers should only install this handler when limiter is non-nil.
func NewHandler(limiter *Server) *Handler {
	return &Handler{limiter: limiter}
}

// Handler applies server-side rate limiting.
type Handler struct {
	limiter *Server
}

// ServeHTTP enforces the configured limiter.
//
// Ignorable paths (health/metrics/etc.) bypass limiting (see `transport/strings.IsIgnorable`).
//
// Behavior:
//   - If `Take` returns an error, it writes an internal server error response.
//   - If `Take` returns a header string, it is added to the response as the "RateLimit" header.
//   - If the request is not allowed, it writes an HTTP 429 response.
//   - Otherwise it calls next.
func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if strings.IsIgnorable(req.URL.Path) {
		next(res, req)
		return
	}

	ctx := req.Context()

	ok, header, err := h.limiter.Take(ctx)
	if err != nil {
		status.WriteError(ctx, res, status.InternalServerError(err))
		return
	}

	res.Header().Add("RateLimit", header)

	if !ok {
		status.WriteError(ctx, res, status.Errorf(http.StatusTooManyRequests, "limiter: too many requests, %s", header))
		return
	}

	next(res, req)
}

// NewClientLimiter constructs an HTTP client-side rate limiter.
//
// If cfg is disabled, it returns (nil, nil) so callers can treat the limiter as not configured.
//
// The returned limiter is backed by `limiter.NewLimiter` and is registered with the provided lifecycle.
// The `keys` map controls how request contexts are turned into limiter keys.
func NewClientLimiter(lc di.Lifecycle, keys KeyMap, cfg *limiter.Config) (*Client, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}

	limiter, err := limiter.NewLimiter(lc, keys, cfg)
	if err != nil {
		return nil, err
	}

	return &Client{limiter}, nil
}

// Client wraps `*limiter.Limiter` for HTTP client integration.
type Client struct {
	*limiter.Limiter
}

// NewRoundTripper constructs an HTTP RoundTripper that enforces rate limiting on outbound requests.
//
// The returned RoundTripper calls `limiter.Take` before delegating to the underlying transport.
// Callers should only install this RoundTripper when limiter is non-nil.
func NewRoundTripper(limiter *Client, rt http.RoundTripper) *RoundTripper {
	return &RoundTripper{limiter: limiter, RoundTripper: rt}
}

// RoundTripper wraps an underlying `http.RoundTripper` and applies client-side rate limiting.
type RoundTripper struct {
	limiter *Client
	http.RoundTripper
}

// RoundTrip enforces the configured limiter.
//
// Behavior:
//   - If `Take` returns an error, RoundTrip returns that error.
//   - If the request is not allowed, it returns an HTTP 429 status error.
//   - Otherwise, it delegates to the underlying RoundTripper.
func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()

	ok, header, err := r.limiter.Take(ctx)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, status.Errorf(http.StatusTooManyRequests, "limiter: too many requests, %s", header)
	}

	return r.RoundTripper.RoundTrip(req)
}
