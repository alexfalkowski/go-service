package meta

import (
	"context"
	"fmt"
	"net/http"

	"github.com/alexfalkowski/go-service/meta"
	m "github.com/alexfalkowski/go-service/transport/meta"
	"github.com/google/uuid"
)

// Handler for meta.
type Handler struct {
	userAgent string
}

// NewHandler for meta.
func NewHandler(userAgent string) *Handler {
	return &Handler{userAgent: userAgent}
}

func (h *Handler) ServeHTTP(resp http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	ctx := req.Context()
	ctx = m.WithUserAgent(ctx, extractUserAgent(ctx, req, h.userAgent))

	requestID := extractRequestID(ctx, req)
	if meta.IsBlank(requestID) {
		requestID = uuid.New()
	}

	ctx = m.WithRequestID(ctx, requestID)

	next(resp, req.WithContext(ctx))
}

// NewRoundTripper for meta.
func NewRoundTripper(userAgent string, hrt http.RoundTripper) *RoundTripper {
	return &RoundTripper{userAgent: userAgent, RoundTripper: hrt}
}

// RoundTripper for meta.
type RoundTripper struct {
	userAgent string
	http.RoundTripper
}

func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()

	userAgent := m.UserAgent(ctx)
	if meta.IsBlank(userAgent) {
		userAgent = meta.Value(r.userAgent)
	}

	req.Header.Set("User-Agent", userAgent.String())
	ctx = m.WithUserAgent(ctx, userAgent)

	requestID := m.RequestID(ctx)
	if meta.IsBlank(requestID) {
		requestID = uuid.New()
	}

	req.Header.Set("Request-ID", requestID.String())
	ctx = m.WithRequestID(ctx, requestID)

	return r.RoundTripper.RoundTrip(req.WithContext(ctx))
}

func extractUserAgent(ctx context.Context, req *http.Request, userAgent string) fmt.Stringer {
	if userAgent := req.Header.Get("User-Agent"); userAgent != "" {
		return meta.Value(userAgent)
	}

	if u := m.UserAgent(ctx); u != nil {
		return u
	}

	return meta.Value(userAgent)
}

func extractRequestID(ctx context.Context, req *http.Request) fmt.Stringer {
	if requestID := req.Header.Get("Request-ID"); requestID != "" {
		return meta.Value(requestID)
	}

	return m.RequestID(ctx)
}
