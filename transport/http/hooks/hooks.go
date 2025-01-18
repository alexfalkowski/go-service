package hooks

import (
	"net/http"
	"strconv"
	"time"

	"github.com/alexfalkowski/go-service/io"
	nh "github.com/alexfalkowski/go-service/net/http"
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
func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	payload, body, err := io.ReadAll(req.Body)
	if err != nil {
		nh.WriteError(req.Context(), res, err, http.StatusBadRequest)

		return
	}

	req.Body = body

	if err := h.hook.Verify(payload, req.Header); err != nil {
		nh.WriteError(req.Context(), res, err, http.StatusBadRequest)

		return
	}

	next(res, req)
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
	payload, body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	req.Body = body

	now := time.Now()
	id := uuid.New().String()

	// Sign does not return an error.
	signature, _ := r.hook.Sign(id, now, payload)

	req.Header.Add(hooks.HeaderWebhookID, id)
	req.Header.Add(hooks.HeaderWebhookSignature, signature)
	req.Header.Add(hooks.HeaderWebhookTimestamp, strconv.FormatInt(now.Unix(), 10))

	return r.RoundTripper.RoundTrip(req)
}
