package limiter

import (
	"net/http"

	"github.com/alexfalkowski/go-service/v2/limiter"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/transport/strings"
)

// Handler for limiter.
func NewHandler(limiter *limiter.Limiter) *Handler {
	return &Handler{limiter: limiter}
}

// Handler for tracer.
type Handler struct {
	limiter *limiter.Limiter
}

// ServeHTTP for limiter.
func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if strings.IsObservable(req.URL.Path) {
		next(res, req)

		return
	}

	ctx := req.Context()

	ok, info, err := h.limiter.Take(ctx)
	if err != nil {
		status.WriteError(ctx, res, status.InternalServerError(err))

		return
	}

	res.Header().Add("RateLimit", info)

	if !ok {
		status.WriteError(ctx, res, status.Errorf(http.StatusTooManyRequests, "limiter: too many requests, %s", info))

		return
	}

	next(res, req)
}
