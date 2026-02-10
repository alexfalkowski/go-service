package limiter

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/limiter"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/transport/strings"
)

// KeyMap is an alias for limiter.KeyMap.
type KeyMap = limiter.KeyMap

// NewServerLimiter returns a server-side rate limiter when enabled.
//
// If cfg is disabled, it returns (nil, nil).
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

// Server wraps limiter.Limiter for HTTP server integration.
type Server struct {
	*limiter.Limiter
}

// NewHandler constructs a server-side rate limiting handler.
func NewHandler(limiter *Server) *Handler {
	return &Handler{limiter: limiter}
}

// Handler applies server-side rate limiting.
type Handler struct {
	limiter *Server
}

// ServeHTTP enforces the limiter and writes HTTP 429 when the rate limit is exceeded.
//
// Requests with ignorable paths bypass limiting.
// When a limit header is returned, it is added to the response as the "RateLimit" header.
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

// NewClientLimiter returns a client-side rate limiter when enabled.
//
// If cfg is disabled, it returns (nil, nil).
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

// Client wraps limiter.Limiter for HTTP client integration.
type Client struct {
	*limiter.Limiter
}

// NewRoundTripper constructs an HTTP RoundTripper that enforces rate limiting on outbound requests.
func NewRoundTripper(limiter *Client, rt http.RoundTripper) *RoundTripper {
	return &RoundTripper{limiter: limiter, RoundTripper: rt}
}

// RoundTripper wraps an underlying http.RoundTripper and applies client-side rate limiting.
type RoundTripper struct {
	limiter *Client
	http.RoundTripper
}

// RoundTrip enforces the limiter and returns HTTP 429 when the rate limit is exceeded.
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
