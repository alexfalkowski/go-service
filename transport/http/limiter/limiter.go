package limiter

import (
	"strconv"

	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/transport/limiter"
)

// KeyMap is an alias for [limiter.KeyMap].
//
// It maps limiter key kinds (for example, "user-agent" or "ip") to functions that derive a rate-limit key
// from the request context.
type KeyMap = limiter.KeyMap

// NewServerLimiter constructs an HTTP server-side rate limiter.
//
// If cfg is disabled, it returns (nil, nil) so callers can treat the limiter as not configured.
//
// The returned limiter is backed by [limiter.NewLimiter] and is registered with the provided lifecycle.
// The `keys` map controls how request contexts are turned into limiter keys (for example, per user-agent).
func NewServerLimiter(lc di.Lifecycle, keys KeyMap, cfg *limiter.Config) (*Server, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}

	rateLimiter, err := limiter.NewLimiter(lc, keys, cfg)
	if err != nil {
		return nil, err
	}

	return &Server{rateLimiter}, nil
}

// Server wraps *[limiter.Limiter] for HTTP server integration.
type Server struct {
	*limiter.Limiter
}

// NewHandler constructs a server-side rate limiting Negroni handler.
//
// Callers should only install this handler when limiter is non-nil.
func NewHandler(routePolicy *http.RoutePolicy, limiter *Server) *Handler {
	return &Handler{routePolicy: routePolicy, limiter: limiter}
}

// Handler applies server-side rate limiting.
type Handler struct {
	limiter     *Server
	routePolicy *http.RoutePolicy
}

// ServeHTTP enforces the configured limiter.
//
// Registered operation paths (health/metrics/etc.) bypass limiting.
//
// Behavior:
//   - If the limiter returns an error, it writes an internal server error response.
//   - It writes RateLimit and RateLimit-Policy headers describing the current decision.
//   - If the request is not allowed, it writes an HTTP 429 response with Retry-After when reset timing is available.
//   - Otherwise it calls next.
func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if h.routePolicy.IsOperation(req) {
		next(res, req)
		return
	}

	ctx := req.Context()

	decision, err := h.limiter.TakeDecision(ctx)
	if err != nil {
		_ = status.WriteError(res, status.InternalServerError(err))
		return
	}

	res.Header().Set("RateLimit", decision.Header())
	res.Header().Set("RateLimit-Policy", decision.PolicyHeader())

	if !decision.Allowed() {
		if resetAfter := decision.ResetAfterSeconds(); resetAfter > 0 {
			res.Header().Set("Retry-After", strconv.FormatUint(resetAfter, 10))
		}

		_ = status.WriteError(res, status.Error(http.StatusTooManyRequests, http.StatusText(http.StatusTooManyRequests)))
		return
	}

	next(res, req)
}

// NewClientLimiter constructs an HTTP client-side rate limiter.
//
// If cfg is disabled, it returns (nil, nil) so callers can treat the limiter as not configured.
//
// The returned limiter is backed by [limiter.NewLimiter] and is registered with the provided lifecycle.
// The `keys` map controls how request contexts are turned into limiter keys.
func NewClientLimiter(lc di.Lifecycle, keys KeyMap, cfg *limiter.Config) (*Client, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}

	rateLimiter, err := limiter.NewLimiter(lc, keys, cfg)
	if err != nil {
		return nil, err
	}

	return &Client{rateLimiter}, nil
}

// Client wraps *[limiter.Limiter] for HTTP client integration.
type Client struct {
	*limiter.Limiter
}

// NewRoundTripper constructs an HTTP RoundTripper that enforces rate limiting on outbound requests.
//
// The returned RoundTripper calls [limiter.Take] before delegating to the underlying transport.
// Callers should only install this RoundTripper when limiter is non-nil.
func NewRoundTripper(limiter *Client, rt http.RoundTripper) *RoundTripper {
	return &RoundTripper{limiter: limiter, RoundTripper: rt}
}

// RoundTripper wraps an underlying [http.RoundTripper] and applies client-side rate limiting.
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
	return http.ClosingRoundTripper(r.roundTrip).RoundTrip(req)
}

func (r *RoundTripper) roundTrip(req *http.Request) (*http.Response, error, bool) {
	ctx := req.Context()

	ok, _, err := r.limiter.Take(ctx)
	if err != nil {
		return nil, err, true
	}

	if !ok {
		return nil, status.LocalError(status.Error(http.StatusTooManyRequests, http.StatusText(http.StatusTooManyRequests))), true
	}

	res, err := r.RoundTripper.RoundTrip(req)
	return res, err, false
}
