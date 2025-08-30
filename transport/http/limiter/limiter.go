package limiter

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/limiter"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/transport/strings"
)

// KeyMap is just an alias for limiter.KeyMap.
type KeyMap = limiter.KeyMap

// NewServerLimiter for http.
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

// Server limiter.
type Server struct {
	*limiter.Limiter
}

// Handler for limiter.
func NewHandler(limiter *Server) *Handler {
	return &Handler{limiter: limiter}
}

// Handler for tracer.
type Handler struct {
	limiter *Server
}

// ServeHTTP for limiter.
func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	p := http.Path(req)
	if strings.IsObservable(p) {
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

// NewClientLimiter for http.
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

// Client limiter.
type Client struct {
	*limiter.Limiter
}

// NewRoundTripper for limiter.
func NewRoundTripper(limiter *Client, rt http.RoundTripper) *RoundTripper {
	return &RoundTripper{limiter: limiter, RoundTripper: rt}
}

// RoundTripper for limiter.
type RoundTripper struct {
	limiter *Client
	http.RoundTripper
}

// RoundTrip for limiter.
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
