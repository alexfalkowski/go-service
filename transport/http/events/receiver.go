package events

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/transport/http/events/hooks"
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
	protocol, _ := events.NewHTTP()

	var handler http.Handler
	handler, _ = events.NewHTTPReceiveHandler(ctx, protocol, receiver)
	handler = hooks.NewHandler(r.hook, handler)

	http.Handle(r.mux, strings.Join(strings.Space, http.MethodPost, path), handler)
}
