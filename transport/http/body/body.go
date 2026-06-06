package body

import (
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/body"
)

// NewHandler constructs request body size limiting middleware.
func NewHandler(limit int64) *Handler {
	return &Handler{limit: limit}
}

// Handler limits request bodies before downstream handlers decode them.
type Handler struct {
	limit int64
}

// ServeHTTP enforces the configured request body limit.
func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	body.NewHandler(next, h.limit).ServeHTTP(res, req)
}
