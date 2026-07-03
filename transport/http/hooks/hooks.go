package hooks

import (
	"strconv"

	"github.com/alexfalkowski/go-service/v2/hooks"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/body"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/time"
	webhooks "github.com/standard-webhooks/standard-webhooks/libraries/go"
)

// NewWebhook constructs a Webhook signer/verifier.
//
// Disabled behavior:
// If hook is nil, NewWebhook returns nil so callers can treat webhook signing/verification as disabled.
//
// Enabled behavior:
// When hook is non-nil, the returned Webhook adds go-service request-body buffering and header conventions
// on top of the underlying Standard Webhooks signer/verifier.
//
// Signing requires a non-nil generator. Passing a nil generator is valid for verification-only use, but
// calling [Webhook.Sign] or using the signing [RoundTripper] with a nil generator will panic.
func NewWebhook(hook *hooks.Hook, generator id.Generator) *Webhook {
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
// The underlying `hook` must be configured with your shared secret(s).
// The `generator` is used to mint webhook ids during signing.
type Webhook struct {
	hook      *hooks.Hook
	generator id.Generator
}

// Sign signs an outbound webhook request.
//
// It reads and buffers the request body, closes the original body, replaces `req.Body` with a fresh
// readable body, and then sets signature headers.
//
// A non-nil generator is required because Sign generates the Webhook-Id used in the signature.
//
// Headers set:
//   - `Webhook-Id`
//   - `Webhook-Signature`
//   - `Webhook-Timestamp`
//
// The signature is computed over the request payload and includes a generated webhook id and timestamp.
// The request body is buffered in memory. Under supported server wiring, inbound request size is already
// capped by the HTTP transport's MaxReceiveSize body limiter before handlers run; outbound requests are
// expected to be bounded by the request construction path.
//
// Disabled behavior:
// If the receiver or underlying Standard Webhooks hook is nil, Sign is a no-op and returns nil.
func (h *Webhook) Sign(req *http.Request) error {
	if h == nil || h.hook == nil {
		return nil
	}

	originalBody := req.Body
	defer body.Close(originalBody)

	payload, bufferedBody, err := body.ReadAll(req)
	if err != nil {
		return err
	}

	req.Body = bufferedBody
	if req.Header == nil {
		req.Header = http.Header{}
	}
	now := time.Now()
	id := h.generator.Generate()
	signature, err := h.hook.Sign(id, now, payload)
	if err != nil {
		return err
	}

	req.Header.Set(webhooks.HeaderWebhookID, id)
	req.Header.Set(webhooks.HeaderWebhookSignature, signature)
	req.Header.Set(webhooks.HeaderWebhookTimestamp, strconv.FormatInt(now.Unix(), 10))

	return nil
}

// Verify verifies the signature headers on an inbound webhook request.
//
// It reads and buffers the request body, closes the original body, replaces `req.Body` with a fresh
// readable body, and then verifies the signature headers using the underlying Standard Webhooks verifier.
//
// Security note:
// This method intentionally does not add a second local size cap. Supported inbound use goes through
// [github.com/alexfalkowski/go-service/v2/transport/http.NewServer], whose body middleware enforces Config.MaxReceiveSize before mux handlers and
// webhook verification run. Callers that bypass the transport server are responsible for applying the same
// request-level cap before invoking Verify.
//
// Replay protection:
// Verify does not persist or reject previously seen Webhook-Id values. Receivers that perform non-idempotent
// work must deduplicate or process idempotently using the Webhook-Id or event identifier.
//
// Disabled behavior:
// If the receiver or underlying Standard Webhooks hook is nil, Verify is a no-op and returns nil.
func (h *Webhook) Verify(req *http.Request) error {
	if h == nil || h.hook == nil {
		return nil
	}

	originalBody := req.Body
	defer body.Close(originalBody)

	payload, bufferedBody, err := body.ReadAll(req)
	if err != nil {
		return err
	}

	req.Body = bufferedBody

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
// Body limits:
// With the standard HTTP transport server, this handler runs after the request-body limiter, so verification
// never sees a body larger than Config.MaxReceiveSize. Direct use outside that server chain must preserve the
// same request-size invariant.
//
// Replay protection:
// Signature verification does not maintain replay state. The wrapped handler is responsible for idempotency
// or deduplication when duplicate valid deliveries would be unsafe.
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
// A non-nil generator must have been supplied to [NewWebhook] for signing use.
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
//
// Redirect behavior:
// If the request is a cross-origin redirect, RoundTrip returns [http.ErrUseLastResponse] without signing
// the redirected request.
func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return http.ClosingRoundTripper(r.roundTrip).RoundTrip(req)
}

func (r *RoundTripper) roundTrip(req *http.Request) (*http.Response, error, bool) {
	if r.hook == nil {
		res, err := r.RoundTripper.RoundTrip(req)
		return res, err, false
	}

	if http.IsCrossOriginRedirect(req) {
		return nil, http.ErrUseLastResponse, true
	}

	clonedReq := req.Clone(req.Context())
	if err := r.hook.Sign(clonedReq); err != nil {
		return nil, err, false
	}

	res, err := r.RoundTripper.RoundTrip(clonedReq)
	return res, err, false
}
