package limiter

import (
	"fmt"
	"net/http"

	"github.com/alexfalkowski/go-service/limiter"
	nh "github.com/alexfalkowski/go-service/net/http"
	"github.com/alexfalkowski/go-service/transport/strings"
	l "github.com/sethvargo/go-limiter"
)

// Handler for limiter.
func NewHandler(limiter l.Store, key limiter.KeyFunc) *Handler {
	return &Handler{limiter: limiter, key: key}
}

// Handler for tracer.
type Handler struct {
	limiter l.Store
	key     limiter.KeyFunc
}

// ServeHTTP for limiter.
//
//nolint:err113
func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if strings.IsObservable(req.URL.Path) {
		next(res, req)

		return
	}

	ctx := req.Context()

	ok, info, err := limiter.Take(ctx, h.limiter, h.key)
	if err != nil {
		nh.WriteError(ctx, res, err, http.StatusInternalServerError)

		return
	}

	res.Header().Add("RateLimit", info)

	if !ok {
		nh.WriteError(ctx, res, fmt.Errorf("limiter: too many requests, %s", info), http.StatusTooManyRequests)

		return
	}

	next(res, req)
}
