package events

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/transport/http/hooks"
	events "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/binding"
	"github.com/cloudevents/sdk-go/v2/client"
	"github.com/cloudevents/sdk-go/v2/protocol"
	transport "github.com/cloudevents/sdk-go/v2/protocol/http"
)

// NewSender constructs a CloudEvents HTTP client that can sign outbound requests with the webhook hook.
//
// The constructed client uses the CloudEvents SDK HTTP transport. Outbound requests are wrapped with the
// webhook signing RoundTripper (see [github.com/alexfalkowski/go-service/v2/transport/http/hooks]) so each request is signed before being sent
// when hook is non-nil. When hook is nil, the signing RoundTripper behaves as a pass-through wrapper and
// requests are sent unsigned.
//
// If no round tripper is configured via [WithSenderRoundTripper], it uses the default HTTP transport.
//
// Note: the provided hook must be configured appropriately (for example with the expected secret) for
// signing to succeed.
func NewSender(hook *hooks.Webhook, opts ...SenderOption) *Sender {
	resolved := options(opts...)
	rt := hooks.NewRoundTripper(hook, resolved.roundTripper)
	httpClient := http.Client{Transport: rt, CheckRedirect: http.SameOriginRedirect, Timeout: resolved.timeout.Duration()}

	sender, _ := events.NewClientHTTP(transport.WithClient(httpClient))
	return &Sender{client: sender}
}

// Sender wraps a CloudEvents client and forces structured HTTP encoding for outbound events.
type Sender struct {
	client client.Client
}

// Send transmits event using structured CloudEvents encoding.
func (s *Sender) Send(ctx context.Context, event events.Event) protocol.Result {
	return s.client.Send(binding.WithForceStructured(ctx), event)
}
