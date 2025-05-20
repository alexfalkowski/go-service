package hooks

import (
	"net/http"

	"github.com/alexfalkowski/go-service/v2/transport/http/hooks"
)

// Handler for hooks.
type Handler struct {
	handler *hooks.Handler

	http.Handler
}

// NewHandler for hooks.
func NewHandler(hook *hooks.Webhook, handler http.Handler) *Handler {
	return &Handler{handler: hooks.NewHandler(hook), Handler: handler}
}

// ServeHTTP for hooks.
func (h *Handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	h.handler.ServeHTTP(resp, req, h.Handler.ServeHTTP)
}
