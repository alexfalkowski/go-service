package limiter

import (
	"fmt"
	"net/http"

	"github.com/alexfalkowski/go-service/limiter"
	"github.com/alexfalkowski/go-service/meta"
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
func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if strings.IsObservable(req.URL.Path) {
		next(res, req)

		return
	}

	ctx := req.Context()

	tokens, remaining, reset, ok, err := h.limiter.Take(ctx, meta.ValueOrBlank(h.key(ctx)))
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)

		return
	}

	res.Header().Add("RateLimit", fmt.Sprintf("limit=%d, remaining=%d, reset=%d", tokens, remaining, reset))

	if !ok {
		http.Error(res, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)

		return
	}

	next(res, req)
}
