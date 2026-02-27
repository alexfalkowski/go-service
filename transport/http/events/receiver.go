package events

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/transport/http/events/hooks"
	events "github.com/cloudevents/sdk-go/v2"
)

// ReceiverFunc is invoked for each received CloudEvent.
//
// The provided context is the request-scoped context associated with the inbound HTTP request.
type ReceiverFunc func(ctx context.Context, e events.Event)

// Receiver registers HTTP handlers that receive CloudEvents.
//
// It wires CloudEvents SDK receive handlers onto a mux and optionally wraps them with webhook verification
// middleware (see `transport/http/events/hooks`).
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
// signatures before allowing events through.
func (r *Receiver) Register(ctx context.Context, path string, receiver ReceiverFunc) {
	protocol, _ := events.NewHTTP()

	var handler http.Handler
	handler, _ = events.NewHTTPReceiveHandler(ctx, protocol, receiver)
	handler = hooks.NewHandler(r.hook, handler)

	http.Handle(r.mux, strings.Join(strings.Space, http.MethodPost, path), handler)
}
