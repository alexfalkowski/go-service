package events

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/events"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/transport/http/events/hooks"
)

// ReceiverFunc is invoked for each received CloudEvent.
//
// The provided context is the request-scoped context associated with the inbound HTTP request.
// Return nil or an ACK result when processing succeeds. Return a non-ACK result, such as a CloudEvents
// HTTP result, when the receiver needs the delivery response to report failure.
type ReceiverFunc = events.ReceiverFunc

// Receiver registers HTTP handlers that receive CloudEvents.
//
// It wires CloudEvents receive handlers onto a mux and optionally wraps them with webhook verification
// middleware (see [github.com/alexfalkowski/go-service/v2/transport/http/events/hooks]).
type Receiver struct {
	mux  *http.ServeMux
	hook *hooks.Webhook
}

// NewReceiver constructs a Receiver that registers handlers on mux.
//
// Registered handlers are wrapped with webhook verification middleware created from hook (if configured).
func NewReceiver(mux *http.ServeMux, hook *hooks.Webhook) *Receiver {
	return &Receiver{mux: mux, hook: hook}
}

// Register registers a CloudEvents HTTP receive handler under path.
//
// The handler is registered for HTTP POST requests and dispatches each decoded event to receiver.
//
// Middleware:
// The receive handler is wrapped with the configured webhook hook handler, which typically verifies request
// signatures before allowing events through. When a webhook hook is configured, the handler supports
// structured HTTP encoding only; binary-mode CloudEvents with ce-* headers are rejected before signature
// verification.
func (r *Receiver) Register(ctx context.Context, path string, receiver ReceiverFunc) {
	handler := events.NewReceiveHandler(ctx, receiver)
	handler = hooks.NewHandler(r.hook, handler)

	http.Handle(r.mux, strings.Join(strings.Space, http.MethodPost, path), handler)
}
