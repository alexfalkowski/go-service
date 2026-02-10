package events

import (
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/transport/http/hooks"
	events "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/client"
	transport "github.com/cloudevents/sdk-go/v2/protocol/http"
)

// SenderOption configures a CloudEvents HTTP sender.
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

// WithSenderRoundTripper configures the underlying HTTP RoundTripper used to send CloudEvents.
func WithSenderRoundTripper(rt http.RoundTripper) SenderOption {
	return senderOptionFunc(func(o *senderOptions) {
		o.roundTripper = rt
	})
}

// NewSender constructs a CloudEvents HTTP client that wraps outbound requests with the webhook hook.
//
// If no round tripper is configured, it uses the default HTTP transport.
func NewSender(hook *hooks.Webhook, opts ...SenderOption) (client.Client, error) {
	os := options(opts...)
	rt := hooks.NewRoundTripper(hook, os.roundTripper)
	return events.NewClientHTTP(transport.WithRoundTripper(rt))
}

func options(opts ...SenderOption) *senderOptions {
	os := &senderOptions{}
	for _, o := range opts {
		o.apply(os)
	}
	if os.roundTripper == nil {
		os.roundTripper = http.Transport(nil)
	}
	return os
}
