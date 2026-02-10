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
func NewWebhook(hook *hooks.Webhook, generator id.Generator) *Webhook {
	return &Webhook{hook: hook, generator: generator}
}

// Webhook signs and verifies webhook requests.
type Webhook struct {
	hook      *hooks.Webhook
	generator id.Generator
}

// Sign reads and buffers the request body, restores req.Body, and then adds signature headers.
//
// The signature is computed over the request payload and includes a generated webhook id and timestamp.
func (h *Webhook) Sign(req *http.Request) error {
	payload, body, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}

	req.Body = body
	now := time.Now()
	id := h.generator.Generate()
	signature, _ := h.hook.Sign(id, now, payload)

	req.Header.Add(hooks.HeaderWebhookID, id)
	req.Header.Add(hooks.HeaderWebhookSignature, signature)
	req.Header.Add(hooks.HeaderWebhookTimestamp, strconv.FormatInt(now.Unix(), 10))

	return nil
}

// Verify reads and buffers the request body, restores req.Body, and verifies the signature headers.
func (h *Webhook) Verify(req *http.Request) error {
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
type Handler struct {
	hook *Webhook
}

// NewHandler constructs a webhook verification handler.
func NewHandler(hook *Webhook) *Handler {
	return &Handler{hook: hook}
}

// ServeHTTP verifies the webhook signature before calling next.
//
// If verification fails, it writes a bad request error response.
func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if err := h.hook.Verify(req); err != nil {
		status.WriteError(req.Context(), res, status.BadRequestError(err))

		return
	}

	next(res, req)
}

// NewRoundTripper constructs an HTTP RoundTripper that signs outbound webhook requests.
func NewRoundTripper(hook *Webhook, rt http.RoundTripper) *RoundTripper {
	return &RoundTripper{hook: hook, RoundTripper: rt}
}

// RoundTripper signs outbound webhook requests before delegating to the underlying RoundTripper.
type RoundTripper struct {
	hook *Webhook
	http.RoundTripper
}

// RoundTrip signs the request and delegates to the underlying RoundTripper.
func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := r.hook.Sign(req); err != nil {
		return nil, err
	}

	return r.RoundTripper.RoundTrip(req)
}
