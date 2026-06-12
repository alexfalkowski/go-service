package events

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/events"
	"github.com/alexfalkowski/go-service/v2/transport/http/hooks"
)

// NewSender constructs a CloudEvents HTTP client that can sign outbound requests with the webhook hook.
//
// The constructed client uses the go-service CloudEvents HTTP transport. Outbound requests are wrapped with the
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

	sender := events.NewClient(httpClient)
	return &Sender{client: sender}
}

// Sender wraps a CloudEvents client and forces structured HTTP encoding for outbound events.
type Sender struct {
	client events.Client
}

// Send transmits event using structured CloudEvents encoding.
func (s *Sender) Send(ctx context.Context, event events.Event) events.Result {
	return events.SendStructured(ctx, s.client, event)
}
