package hooks

import (
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/transport/http/hooks"
)

// Webhook is an alias for hooks.Webhook.
type Webhook = hooks.Webhook

// NewHandler constructs an HTTP handler that applies the webhook hook handler before invoking handler.
func NewHandler(hook *Webhook, handler http.Handler) *Handler {
	return &Handler{handler: hooks.NewHandler(hook), Handler: handler}
}

// Handler wraps an http.Handler and applies a webhook hook handler.
type Handler struct {
	handler *hooks.Handler
	http.Handler
}

// ServeHTTP applies the webhook hook handler and then delegates to the wrapped handler.
func (h *Handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	h.handler.ServeHTTP(resp, req, h.Handler.ServeHTTP)
}
