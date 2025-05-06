package http

import (
	"context"
	"net/http"

	"github.com/alexfalkowski/go-service/strings"
	hh "github.com/alexfalkowski/go-service/transport/events/http/hooks"
	"github.com/alexfalkowski/go-service/transport/http/hooks"
	events "github.com/cloudevents/sdk-go/v2"
)

// ReceiverFunc for HTTP.
type ReceiverFunc func(ctx context.Context, e events.Event)

// Receiver for HTTP.
type Receiver struct {
	mux  *http.ServeMux
	hook *hooks.Webhook
}

// NewReceiver for HTTP.
func NewReceiver(mux *http.ServeMux, hook *hooks.Webhook) *Receiver {
	return &Receiver{mux: mux, hook: hook}
}

// Register a fn under path.
func (r *Receiver) Register(ctx context.Context, path string, receiver ReceiverFunc) {
	// Error is only returned with options.
	protocol, _ := events.NewHTTP()

	var handler http.Handler

	// Error is only returned when an incorrect signature of a function is used (it uses reflection).
	handler, _ = events.NewHTTPReceiveHandler(ctx, protocol, receiver)
	handler = hh.NewHandler(r.hook, handler)

	r.mux.Handle(strings.Join(" ", http.MethodPost, path), handler)
}
