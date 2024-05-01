package http

import (
	"context"
	"net/http"

	h "github.com/alexfalkowski/go-service/transport/events/http/hooks"
	events "github.com/cloudevents/sdk-go/v2"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	hooks "github.com/standard-webhooks/standard-webhooks/libraries/go"
)

// ReceiverFunc for HTTP.
type ReceiverFunc func(ctx context.Context, e events.Event)

// Receiver for HTTP.
type Receiver struct {
	mux  *runtime.ServeMux
	hook *hooks.Webhook
}

// NewReceiver for HTTP.
func NewReceiver(mux *runtime.ServeMux, hook *hooks.Webhook) *Receiver {
	return &Receiver{mux: mux, hook: hook}
}

// Register a fn under path.
func (r *Receiver) Register(ctx context.Context, path string, fn ReceiverFunc) error {
	// Error is only returned with options.
	p, _ := events.NewHTTP()

	var (
		handler http.Handler
		err     error
	)

	handler, err = events.NewHTTPReceiveHandler(ctx, p, fn)
	if err != nil {
		return err
	}

	handler = h.NewHandler(r.hook, handler)

	return r.mux.HandlePath("POST", path, func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		handler.ServeHTTP(w, r)
	})
}
