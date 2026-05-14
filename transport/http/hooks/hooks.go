package hooks

import (
	"strconv"

	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/time"
	hooks "github.com/standard-webhooks/standard-webhooks/libraries/go"
)

// NewWebhook constructs a Webhook signer/verifier.
//
// Disabled behavior:
// If hook is nil, NewWebhook returns nil so callers can treat webhook signing/verification as disabled.
//
// Enabled behavior:
// When hook is non-nil, the returned Webhook adds go-service request-body buffering and header conventions
// on top of the underlying Standard Webhooks signer/verifier.
func NewWebhook(hook *hooks.Webhook, generator id.Generator) *Webhook {
	if hook == nil {
		return nil
	}

	return &Webhook{hook: hook, generator: generator}
}

// Webhook signs and verifies webhook requests using the Standard Webhooks protocol.
//
// It is a thin wrapper around `standard-webhooks` that adds go-service conventions:
//
//   - request body buffering with restoration of `req.Body`
//   - consistent header setting for webhook id, signature, and timestamp
//
// The underlying `hook` must be configured with your shared secret(s) as required by the Standard Webhooks library.
// The `generator` is used to mint webhook ids during signing.
type Webhook struct {
	hook      *hooks.Webhook
	generator id.Generator
}

// Sign signs an outbound webhook request.
//
// It reads and buffers the request body, restores `req.Body`, and then sets signature headers.
//
// Headers set:
//   - `Webhook-Id`
//   - `Webhook-Signature`
//   - `Webhook-Timestamp`
//
// The signature is computed over the request payload and includes a generated webhook id and timestamp.
// Callers should ensure the request body is readable (and reasonably bounded) since it is buffered in memory.
//
// Disabled behavior:
// If the receiver or underlying Standard Webhooks hook is nil, Sign is a no-op and returns nil.
func (h *Webhook) Sign(req *http.Request) error {
	if h == nil || h.hook == nil {
		return nil
	}

	payload, body, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}

	req.Body = body
	now := time.Now()
	id := h.generator.Generate()
	signature, _ := h.hook.Sign(id, now, payload)

	req.Header.Set(hooks.HeaderWebhookID, id)
	req.Header.Set(hooks.HeaderWebhookSignature, signature)
	req.Header.Set(hooks.HeaderWebhookTimestamp, strconv.FormatInt(now.Unix(), 10))

	return nil
}

// Verify verifies the signature headers on an inbound webhook request.
//
// It reads and buffers the request body, restores `req.Body`, and then verifies the signature headers
// using the underlying Standard Webhooks verifier.
//
// Callers should ensure the request body is readable (and reasonably bounded) since it is buffered in memory.
//
// Disabled behavior:
// If the receiver or underlying Standard Webhooks hook is nil, Verify is a no-op and returns nil.
func (h *Webhook) Verify(req *http.Request) error {
	if h == nil || h.hook == nil {
		return nil
	}

	payload, body, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}

	req.Body = body

	if err := h.hook.Verify(payload, req.Header); err != nil {
		return err
	}

	return nil
}

// Handler verifies webhook signatures on inbound requests.
//
// When webhook support is disabled (hook is nil), Handler becomes a pass-through middleware.
type Handler struct {
	hook *Webhook
}

// NewHandler constructs webhook verification middleware.
//
// If hook is nil, the returned handler behaves as pass-through middleware and simply calls next.
func NewHandler(hook *Webhook) *Handler {
	return &Handler{hook: hook}
}

// ServeHTTP verifies the webhook signature before calling next.
//
// If verification fails, it writes an HTTP 400 error response and does not call next.
//
// Disabled behavior:
// If the handler or its hook is nil, ServeHTTP behaves as a pass-through and immediately calls next.
func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if h == nil || h.hook == nil {
		next(res, req)
		return
	}

	if err := h.hook.Verify(req); err != nil {
		_ = status.WriteError(res, status.BadRequestError(err))

		return
	}

	next(res, req)
}

// NewRoundTripper constructs an HTTP RoundTripper that signs outbound webhook requests.
//
// The returned RoundTripper signs each request by buffering the request body, restoring `req.Body`, and
// attaching the Standard Webhooks signature headers before delegating to the underlying transport.
//
// If hook is nil, the returned RoundTripper behaves as a pass-through wrapper.
func NewRoundTripper(hook *Webhook, rt http.RoundTripper) *RoundTripper {
	return &RoundTripper{hook: hook, RoundTripper: rt}
}

// RoundTripper signs outbound webhook requests before delegating to the underlying RoundTripper.
//
// When webhook support is disabled (hook is nil), RoundTripper becomes a pass-through wrapper.
type RoundTripper struct {
	hook *Webhook
	http.RoundTripper
}

// RoundTrip signs the request and delegates to the underlying RoundTripper.
//
// If signing fails (for example, due to an unreadable body), RoundTrip returns the signing error.
//
// Disabled behavior:
// If the configured hook is nil, RoundTrip delegates directly to the underlying RoundTripper without
// mutating the request.
func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if r.hook == nil {
		return r.RoundTripper.RoundTrip(req)
	}

	cloned := req.Clone(req.Context())
	body := cloned.Body
	if err := r.hook.Sign(cloned); err != nil {
		closeBody(body)

		return nil, err
	}
	closeBody(body)

	return r.RoundTripper.RoundTrip(cloned)
}

func closeBody(body io.ReadCloser) {
	if body != nil && body != http.NoBody {
		_ = body.Close()
	}
}
