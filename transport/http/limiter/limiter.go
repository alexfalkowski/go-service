package limiter

import (
	"net/http"

	"github.com/alexfalkowski/go-service/limiter"
	"github.com/alexfalkowski/go-service/net/http/status"
	"github.com/alexfalkowski/go-service/transport/strings"
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
		err := status.FromError(http.StatusInternalServerError, err)
		status.WriteError(res, err)

		return
	}

	res.Header().Add("RateLimit", info)

	if !ok {
		err := status.Errorf(http.StatusTooManyRequests, "limiter: too many requests, %s", info)
		status.WriteError(res, err)

		return
	}

	next(res, req)
}
