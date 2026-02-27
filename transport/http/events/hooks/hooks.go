package hooks

import (
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/transport/http/hooks"
)

// Webhook is an alias for `transport/http/hooks.Webhook`.
//
// It is re-exported from this package to provide a CloudEvents-focused import path while still using the
// shared webhook implementation.
type Webhook = hooks.Webhook

// NewHandler constructs an HTTP handler that verifies webhook signatures before invoking handler.
//
// The returned handler wraps handler and applies the webhook verification middleware first. If signature
// verification fails, the request is rejected and handler is not invoked.
func NewHandler(hook *Webhook, handler http.Handler) *Handler {
	return &Handler{handler: hooks.NewHandler(hook), Handler: handler}
}

// Handler wraps an http.Handler and applies webhook verification middleware.
type Handler struct {
	handler *hooks.Handler
	http.Handler
}

// ServeHTTP verifies the webhook signature and then delegates to the wrapped handler.
//
// If verification fails, the underlying middleware writes an HTTP error response and does not call the wrapped
// handler.
func (h *Handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	h.handler.ServeHTTP(resp, req, h.Handler.ServeHTTP)
}
