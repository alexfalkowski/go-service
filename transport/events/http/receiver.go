package http

import (
	"context"
	"net/http"

	h "github.com/alexfalkowski/go-service/transport/events/http/hooks"
	events "github.com/cloudevents/sdk-go/v2"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	hooks "github.com/standard-webhooks/standard-webhooks/libraries/go"
)

// ReceiverOption for HTTP.
type ReceiverOption interface{ apply(opts *receiverOptions) }

type receiverOptions struct {
	hook *hooks.Webhook
}

type receiverOptionFunc func(*receiverOptions)

func (f receiverOptionFunc) apply(o *receiverOptions) { f(o) }

// WithReceiverHook for HTTP.
func WithReceiverHook(hook *hooks.Webhook) ReceiverOption {
	return receiverOptionFunc(func(o *receiverOptions) {
		o.hook = hook
	})
}

// Receiver for HTTP.
type Receiver func(ctx context.Context, e events.Event)

// RegisterReceiver for HTTP.
func RegisterReceiver(ctx context.Context, mux *runtime.ServeMux, path string, recv Receiver, opts ...ReceiverOption) error {
	os := &receiverOptions{}
	for _, o := range opts {
		o.apply(os)
	}

	p, err := events.NewHTTP()
	if err != nil {
		return err
	}

	var handler http.Handler

	handler, err = events.NewHTTPReceiveHandler(ctx, p, recv)
	if err != nil {
		return err
	}

	if os.hook != nil {
		handler = h.NewHandler(os.hook, handler)
	}

	return mux.HandlePath("POST", path, func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		handler.ServeHTTP(w, r)
	})
}
