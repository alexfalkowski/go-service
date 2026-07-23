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
// when hook is non-nil. Each Send uses one generated Webhook-Id so retries retain the same id while refreshing
// their timestamp and signature. When hook is nil, the signing RoundTripper behaves as a pass-through wrapper
// and requests are sent unsigned.
//
// If no round tripper is configured via [WithSenderRoundTripper], it uses the default HTTP transport.
//
// The sender uses structured CloudEvents HTTP encoding by default. Use [WithSenderEncoding] to select binary
// encoding for outbound integrations that require it.
//
// The sender follows only same-origin redirects. Cross-origin redirects are returned to the caller instead
// of being followed so webhook signatures are not minted for redirected origins.
//
// Note: the provided hook must be configured for signing use, including the expected secret and a non-nil
// id generator passed to [github.com/alexfalkowski/go-service/v2/transport/http/hooks.NewWebhook].
func NewSender(hook *hooks.Webhook, opts ...SenderOption) *Sender {
	resolved := options(opts...)
	rt := hooks.NewRoundTripper(hook, resolved.roundTripper)
	httpClient := http.NewClient(rt, resolved.timeout)
	httpClient.CheckRedirect = http.SameOriginRedirect

	sender := events.NewClient(*httpClient)
	return &Sender{client: sender, encoding: resolved.encoding, hook: hook}
}

// Sender wraps a CloudEvents client and sends outbound events using its configured HTTP encoding.
type Sender struct {
	client   events.Client
	hook     *hooks.Webhook
	encoding SenderEncoding
}

// Send transmits event using the configured CloudEvents HTTP encoding.
func (s *Sender) Send(ctx context.Context, event events.Event) events.Result {
	if s.hook != nil {
		ctx = hooks.WithWebhookID(ctx, s.hook.GenerateID())
	}

	if s.encoding == SenderEncodingBinary {
		return events.SendBinary(ctx, s.client, event)
	}

	return events.SendStructured(ctx, s.client, event)
}
