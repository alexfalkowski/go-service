package http

import (
	"net/http"

	sh "github.com/alexfalkowski/go-service/transport/http"
	h "github.com/alexfalkowski/go-service/transport/http/hooks"
	events "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/client"
	eh "github.com/cloudevents/sdk-go/v2/protocol/http"
	hooks "github.com/standard-webhooks/standard-webhooks/libraries/go"
)

// SenderOption for HTTP.
type SenderOption interface{ apply(opts *senderOptions) }

type senderOptions struct {
	roundTripper http.RoundTripper
	hook         *hooks.Webhook
}

type senderOptionFunc func(*senderOptions)

func (f senderOptionFunc) apply(o *senderOptions) { f(o) }

// WithSenderRoundTripper for HTTP.
func WithSenderRoundTripper(rt http.RoundTripper) SenderOption {
	return senderOptionFunc(func(o *senderOptions) {
		o.roundTripper = rt
	})
}

// WithSenderHook for HTTP.
func WithSenderHook(hook *hooks.Webhook) SenderOption {
	return senderOptionFunc(func(o *senderOptions) {
		o.hook = hook
	})
}

// NewSender for HTTP.
func NewSender(opts ...SenderOption) (client.Client, error) {
	os := &senderOptions{}
	for _, o := range opts {
		o.apply(os)
	}

	rt := os.roundTripper
	if rt == nil {
		rt = sh.Transport()
	}

	if os.hook != nil {
		rt = h.NewRoundTripper(os.hook, rt)
	}

	return events.NewClientHTTP(eh.WithRoundTripper(rt))
}
