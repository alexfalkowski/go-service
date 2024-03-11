package hooks

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/alexfalkowski/go-service/meta"
	"github.com/google/uuid"
	hooks "github.com/standard-webhooks/standard-webhooks/libraries/go"
)

// Handler for hooks.
type Handler struct {
	hook *hooks.Webhook

	http.Handler
}

// NewHandler for hooks.
func NewHandler(hook *hooks.Webhook, handler http.Handler) *Handler {
	return &Handler{hook: hook, Handler: handler}
}

// ServeHTTP for hooks.
func (h *Handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	payload, err := io.ReadAll(req.Body)
	if err != nil {
		meta.WithAttribute(ctx, "hooksError", err.Error())

		return
	}

	req.Body = io.NopCloser(bytes.NewReader(payload))

	if err := h.hook.Verify(payload, req.Header); err != nil {
		meta.WithAttribute(ctx, "hooksError", err.Error())

		return
	}

	h.Handler.ServeHTTP(resp, req)
}

// NewRoundTripper for hooks.
func NewRoundTripper(hook *hooks.Webhook, hrt http.RoundTripper) *RoundTripper {
	return &RoundTripper{hook: hook, RoundTripper: hrt}
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
