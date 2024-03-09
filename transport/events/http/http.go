package http

import (
	"context"
	"net/http"

	sh "github.com/alexfalkowski/go-service/transport/http"
	events "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/client"
	eh "github.com/cloudevents/sdk-go/v2/protocol/http"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

// Option for HTTP.
type Option interface{ apply(opts *options) }

type options struct {
	roundTripper http.RoundTripper
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) { f(o) }

// WithRoundTripper for HTTP.
func WithRoundTripper(rt http.RoundTripper) Option {
	return optionFunc(func(o *options) {
		o.roundTripper = rt
	})
}

// NewSender for HTTP.
func NewSender(opts ...Option) (client.Client, error) {
	os := &options{}
	for _, o := range opts {
		o.apply(os)
	}

	rt := os.roundTripper
	if rt == nil {
		rt = sh.Transport()
	}

	c, err := events.NewClientHTTP(eh.WithRoundTripper(rt))
	if err != nil {
		return nil, err
	}

	return c, err
}

// Receiver for HTTP.
type Receiver func(ctx context.Context, e events.Event)

// RegisterReceiver for HTTP.
func RegisterReceiver(ctx context.Context, path string, mux *runtime.ServeMux, recv Receiver) error {
	p, err := events.NewHTTP()
	if err != nil {
		return err
	}

	h, err := events.NewHTTPReceiveHandler(ctx, p, recv)
	if err != nil {
		return err
	}

	return mux.HandlePath("POST", path, func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		h.ServeHTTP(w, r)
	})
}
