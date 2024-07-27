package limiter

import (
	"fmt"
	"net/http"
	"time"

	"github.com/alexfalkowski/go-service/limiter"
	"github.com/alexfalkowski/go-service/meta"
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
func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if strings.IsObservable(req.URL.Path) {
		next(res, req)

		return
	}

	ctx := req.Context()

	tokens, remaining, reset, ok, err := h.limiter.Take(ctx, meta.ValueOrBlank(h.key(ctx)))
	if err != nil {
		nh.WriteError(req.Context(), res, err, http.StatusInternalServerError)

		return
	}

	r := time.Until(time.Unix(0, int64(reset)))
	v := fmt.Sprintf("limit=%d, remaining=%d, reset=%s", tokens, remaining, r)

	res.Header().Add("RateLimit", v)

	if !ok {
		nh.WriteError(ctx, res, fmt.Errorf("limiter: %s", v), http.StatusTooManyRequests)

		return
	}

	next(res, req)
}
