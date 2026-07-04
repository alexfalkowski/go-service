package events

import (
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/time"
)

// SenderEncoding selects the CloudEvents HTTP encoding used by a sender.
type SenderEncoding uint8

// SenderEncodingStructured sends CloudEvents using structured HTTP encoding.
const SenderEncodingStructured SenderEncoding = 0

// SenderEncodingBinary sends CloudEvents using binary HTTP encoding.
const SenderEncodingBinary SenderEncoding = 1

// SenderOption configures a CloudEvents HTTP sender.
//
// Sender options control how the CloudEvents client sends events, including the underlying HTTP transport and
// CloudEvents HTTP encoding.
type SenderOption interface {
	apply(opts *senderOptions)
}

type senderOptions struct {
	roundTripper http.RoundTripper
	timeout      time.Duration
	encoding     SenderEncoding
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

// WithSenderTimeout configures the HTTP client timeout used when sending CloudEvents.
//
// A non-positive timeout uses [time.DefaultTimeout].
func WithSenderTimeout(timeout time.Duration) SenderOption {
	return senderOptionFunc(func(o *senderOptions) {
		o.timeout = timeout
	})
}

// WithSenderEncoding configures the CloudEvents HTTP encoding used when sending events.
//
// The default is [SenderEncodingStructured].
// Webhook-protected go-service receivers require structured encoding and reject binary-mode CloudEvents.
func WithSenderEncoding(encoding SenderEncoding) SenderOption {
	return senderOptionFunc(func(o *senderOptions) {
		o.encoding = encoding
	})
}

func options(opts ...SenderOption) *senderOptions {
	resolved := &senderOptions{}
	for _, o := range opts {
		o.apply(resolved)
	}

	if resolved.roundTripper == nil {
		resolved.roundTripper = http.Transport(nil)
	}

	if resolved.timeout <= 0 {
		resolved.timeout = time.DefaultTimeout
	}

	return resolved
}
