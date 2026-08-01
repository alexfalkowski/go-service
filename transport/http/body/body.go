package body

import (
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/body"
)

// NewHandler constructs request body size limiting middleware that never buffers routes marked streaming in routePolicy,
// and buffers every other route.
//
// The limit is expressed in bytes and is passed to the underlying
// [github.com/alexfalkowski/go-service/v2/net/http/body.NewHandler] and
// [github.com/alexfalkowski/go-service/v2/net/http/body.NewLazyHandler] wrappers. Unlike the buffered path,
// the streaming branch does not enforce limit as a cumulative body-size ceiling; see
// [github.com/alexfalkowski/go-service/v2/net/http/body.NewLazyHandler].
func NewHandler(routePolicy *http.RoutePolicy, limit int64) *Handler {
	return &Handler{routePolicy: routePolicy, limit: limit}
}

// Handler limits request bodies before downstream handlers decode them.
//
// Accepted non-empty bodies are buffered and replaced with a fresh readable body for downstream
// handlers, unless routePolicy marks the request's route as streaming, in which case the body is
// never buffered and no cumulative limit is enforced as it is read (see
// [github.com/alexfalkowski/go-service/v2/net/http/body.NewLazyHandler]). Empty bodies are passed
// through unchanged.
type Handler struct {
	routePolicy *http.RoutePolicy
	limit       int64
}

// ServeHTTP enforces the configured request body limit.
//
// If the request body exceeds the configured limit, ServeHTTP writes the
// underlying HTTP max-bytes error response and does not call next. If the body
// cannot be read, ServeHTTP writes a bad-request response and does not call
// next. Otherwise, it delegates to next with a readable request body.
//
// When routePolicy marks req's route as streaming, ServeHTTP never buffers the body (see
// [github.com/alexfalkowski/go-service/v2/net/http/body.NewLazyHandler]): it still rejects a request
// whose declared Content-Length exceeds the limit before next runs, but otherwise delegates to next
// without enforcing any cumulative limit as next reads the body. A streaming route's per-value cap, if
// any, is applied by next itself (see
// [github.com/alexfalkowski/go-service/v2/net/http/content/stream.RequestStream.Recv]), not by ServeHTTP.
func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if h.routePolicy != nil && h.routePolicy.IsStreaming(req) {
		body.NewLazyHandler(next, h.limit).ServeHTTP(res, req)
		return
	}

	body.NewHandler(next, h.limit).ServeHTTP(res, req)
}
