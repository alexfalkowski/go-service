package hooks

import (
	"net/http"

	h "github.com/alexfalkowski/go-service/transport/http/hooks"
	hooks "github.com/standard-webhooks/standard-webhooks/libraries/go"
)

// Handler for hooks.
type Handler struct {
	handler *h.Handler

	http.Handler
}

// NewHandler for hooks.
func NewHandler(hook *hooks.Webhook, handler http.Handler) *Handler {
	return &Handler{handler: h.NewHandler(hook), Handler: handler}
}

// ServeHTTP for hooks.
func (h *Handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	h.handler.ServeHTTP(resp, req, h.Handler.ServeHTTP)
}
