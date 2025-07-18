package limiter

import (
	"github.com/alexfalkowski/go-service/v2/limiter"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/transport/strings"
)

// Limiter is just an alias for limiter.Limiter.
type Limiter = limiter.Limiter

// Handler for limiter.
func NewHandler(limiter *Limiter) *Handler {
	return &Handler{limiter: limiter}
}

// Handler for tracer.
type Handler struct {
	limiter *Limiter
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

// NewRoundTripper for limiter.
func NewRoundTripper(limiter *Limiter, rt http.RoundTripper) *RoundTripper {
	return &RoundTripper{limiter: limiter, RoundTripper: rt}
}

// RoundTripper for limiter.
type RoundTripper struct {
	limiter *Limiter
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
