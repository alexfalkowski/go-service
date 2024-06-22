package limiter

import (
	"net/http"
	"strconv"

	"github.com/alexfalkowski/go-service/limiter"
	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/transport/strings"
	l "github.com/ulule/limiter/v3"
)

// Handler for limiter.
func NewHandler(limiter *l.Limiter, key limiter.KeyFunc) *Handler {
	return &Handler{limiter: limiter, key: key}
}

// Handler for tracer.
type Handler struct {
	limiter *l.Limiter
	key     limiter.KeyFunc
}

// ServeHTTP for limiter.
func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if strings.IsObservable(req.URL.Path) {
		next(res, req)

		return
	}

	ctx := req.Context()

	// Memory stores do not return error.
	context, _ := h.limiter.Get(ctx, meta.ValueOrBlank(h.key(ctx)))

	if context.Reached {
		res.Header().Add("X-Rate-Limit-Limit", strconv.FormatInt(context.Limit, 10))
		res.WriteHeader(http.StatusTooManyRequests)

		return
	}

	next(res, req)
}
