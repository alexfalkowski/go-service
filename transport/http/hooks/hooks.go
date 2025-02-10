package hooks

import (
	"net/http"
	"strconv"

	"github.com/alexfalkowski/go-service/id"
	"github.com/alexfalkowski/go-service/io"
	nh "github.com/alexfalkowski/go-service/net/http"
	"github.com/alexfalkowski/go-service/time"
	hooks "github.com/standard-webhooks/standard-webhooks/libraries/go"
)

// Webhook provides a simple facade that signs and verifies the payload.
type Webhook struct {
	hook *hooks.Webhook
	gen  id.Generator
}

// NewWebhook for http.
func NewWebhook(hook *hooks.Webhook, gen id.Generator) *Webhook {
	return &Webhook{hook: hook, gen: gen}
}

// Sign the webhook.
func (h *Webhook) Sign(req *http.Request) error {
	payload, body, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}

	req.Body = body

	now := time.Now()
	id := h.gen.Generate()

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
		nh.WriteError(req.Context(), res, err, http.StatusBadRequest)

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
