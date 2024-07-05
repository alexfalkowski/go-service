package http

import (
	"net/http"

	nh "github.com/alexfalkowski/go-service/net/http"
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
}

type senderOptionFunc func(*senderOptions)

func (f senderOptionFunc) apply(o *senderOptions) { f(o) }

// WithSenderRoundTripper for HTTP.
func WithSenderRoundTripper(rt http.RoundTripper) SenderOption {
	return senderOptionFunc(func(o *senderOptions) {
		o.roundTripper = rt
	})
}

// NewSender for HTTP.
func NewSender(hook *hooks.Webhook, opts ...SenderOption) (client.Client, error) {
	os := &senderOptions{roundTripper: nh.Transport(nil)}
	for _, o := range opts {
		o.apply(os)
	}

	rt := os.roundTripper
	rt = h.NewRoundTripper(hook, rt)

	return events.NewClientHTTP(eh.WithRoundTripper(rt))
}
