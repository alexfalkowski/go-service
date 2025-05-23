package events

import (
	"net/http"

	nh "github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/transport/http/hooks"
	events "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/client"
	eh "github.com/cloudevents/sdk-go/v2/protocol/http"
)

// SenderOption for HTTP.
type SenderOption interface {
	apply(opts *senderOptions)
}

type senderOptions struct {
	roundTripper http.RoundTripper
}

type senderOptionFunc func(*senderOptions)

func (f senderOptionFunc) apply(o *senderOptions) {
	f(o)
}

// WithSenderRoundTripper for HTTP.
func WithSenderRoundTripper(rt http.RoundTripper) SenderOption {
	return senderOptionFunc(func(o *senderOptions) {
		o.roundTripper = rt
	})
}

// NewSender for HTTP.
func NewSender(hook *hooks.Webhook, opts ...SenderOption) (client.Client, error) {
	os := options(opts...)

	rt := os.roundTripper
	rt = hooks.NewRoundTripper(hook, rt)

	return events.NewClientHTTP(eh.WithRoundTripper(rt))
}

func options(opts ...SenderOption) *senderOptions {
	os := &senderOptions{}
	for _, o := range opts {
		o.apply(os)
	}

	if os.roundTripper == nil {
		os.roundTripper = nh.Transport(nil)
	}

	return os
}
