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

// NewWebhook for http.
func NewWebhook(hook *hooks.Webhook, generator id.Generator) *Webhook {
	return &Webhook{hook: hook, generator: generator}
}

// Webhook provides a simple facade that signs and verifies the payload.
type Webhook struct {
	hook      *hooks.Webhook
	generator id.Generator
}

// Sign the webhook.
func (h *Webhook) Sign(req *http.Request) error {
	payload, body, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}

	req.Body = body

	now := time.Now()
	id := h.generator.Generate()

	// Sign does not return an error.
	signature, _ := h.hook.Sign(id, now, payload)

	req.Header.Add(hooks.HeaderWebhookID, id)
	req.Header.Add(hooks.HeaderWebhookSignature, signature)
	req.Header.Add(hooks.HeaderWebhookTimestamp, strconv.FormatInt(now.Unix(), 10))

	return nil
}

// Verify the webhook.
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

// Handler for hooks.
type Handler struct {
	hook *Webhook
}

// NewHandler for hooks.
func NewHandler(hook *Webhook) *Handler {
	return &Handler{hook: hook}
}

// ServeHTTP for hooks.
func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if err := h.hook.Verify(req); err != nil {
		status.WriteError(req.Context(), res, status.BadRequestError(err))

		return
	}

	next(res, req)
}

// NewRoundTripper for hooks.
func NewRoundTripper(hook *Webhook, rt http.RoundTripper) *RoundTripper {
	return &RoundTripper{hook: hook, RoundTripper: rt}
}

// RoundTripper for hooks.
type RoundTripper struct {
	hook *Webhook
	http.RoundTripper
}

// RoundTrip for hooks.
func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := r.hook.Sign(req); err != nil {
		return nil, err
	}

	return r.RoundTripper.RoundTrip(req)
}
