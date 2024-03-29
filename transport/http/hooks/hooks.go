package hooks

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	hooks "github.com/standard-webhooks/standard-webhooks/libraries/go"
)

// Handler for hooks.
type Handler struct {
	hook *hooks.Webhook
}

// NewHandler for hooks.
func NewHandler(hook *hooks.Webhook) *Handler {
	return &Handler{hook: hook}
}

// ServeHTTP for hooks.
func (h *Handler) ServeHTTP(resp http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	payload, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)

		return
	}

	req.Body = io.NopCloser(bytes.NewReader(payload))

	if err := h.hook.Verify(payload, req.Header); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)

		return
	}

	next(resp, req)
}

// NewRoundTripper for hooks.
func NewRoundTripper(hook *hooks.Webhook, rt http.RoundTripper) *RoundTripper {
	return &RoundTripper{hook: hook, RoundTripper: rt}
}

// RoundTripper for hooks.
type RoundTripper struct {
	hook *hooks.Webhook

	http.RoundTripper
}

// RoundTrip for hooks.
func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	payload, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	req.Body = io.NopCloser(bytes.NewReader(payload))

	ts := time.Now()
	id := uuid.New().String()

	signature, err := r.hook.Sign(id, ts, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Add(hooks.HeaderWebhookID, id)
	req.Header.Add(hooks.HeaderWebhookSignature, signature)
	req.Header.Add(hooks.HeaderWebhookTimestamp, strconv.FormatInt(ts.Unix(), 10))

	return r.RoundTripper.RoundTrip(req)
}
