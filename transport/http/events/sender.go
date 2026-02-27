package events

import (
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/transport/http/hooks"
	events "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/client"
	transport "github.com/cloudevents/sdk-go/v2/protocol/http"
)

// SenderOption configures a CloudEvents HTTP sender.
//
// Sender options control how the CloudEvents client is constructed, primarily by selecting the underlying
// HTTP transport used to send events.
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
//
// This is an escape hatch for providing a custom transport (for example, one that is instrumented,
// uses a custom proxy, or is a test double).
func WithSenderRoundTripper(rt http.RoundTripper) SenderOption {
	return senderOptionFunc(func(o *senderOptions) {
		o.roundTripper = rt
	})
}

// NewSender constructs a CloudEvents HTTP client that signs outbound requests with the webhook hook.
//
// The constructed client uses the CloudEvents SDK HTTP transport. Outbound requests are wrapped with the
// webhook signing RoundTripper (see `transport/http/hooks`) so each request is signed before being sent.
//
// If no round tripper is configured via `WithSenderRoundTripper`, it uses the default HTTP transport.
//
// Note: the provided hook must be configured appropriately (for example with the expected secret) for
// signing to succeed.
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
