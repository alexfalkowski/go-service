package body

import (
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/body"
)

// NewHandler constructs request body size limiting middleware.
//
// The limit is expressed in bytes and is passed to the underlying
// [github.com/alexfalkowski/go-service/v2/net/http/body.NewHandler] wrapper.
func NewHandler(limit int64) *Handler {
	return &Handler{limit: limit}
}

// Handler limits request bodies before downstream handlers decode them.
//
// Accepted non-empty bodies are buffered and replaced with a fresh readable
// body for downstream handlers. Empty bodies are passed through unchanged.
type Handler struct {
	limit int64
}

// ServeHTTP enforces the configured request body limit.
//
// If the request body exceeds the configured limit, ServeHTTP writes the
// underlying HTTP max-bytes error response and does not call next. If the body
// cannot be read, ServeHTTP writes a bad-request response and does not call
// next. Otherwise, it delegates to next with a readable request body.
func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	body.NewHandler(next, h.limit).ServeHTTP(res, req)
}
