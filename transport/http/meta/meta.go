package meta

import (
	"context"
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
	ctx = m.WithRequestID(ctx, extractRequestID(ctx, req))

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

	userAgent := extractUserAgent(ctx, req, r.userAgent)

	req.Header.Set("User-Agent", userAgent.Value())
	ctx = m.WithUserAgent(ctx, userAgent)

	requestID := extractRequestID(ctx, req)

	req.Header.Set("Request-ID", requestID.Value())
	ctx = m.WithRequestID(ctx, requestID)

	return r.RoundTripper.RoundTrip(req.WithContext(ctx))
}

func extractUserAgent(ctx context.Context, req *http.Request, userAgent string) meta.Valuer {
	if ua := req.Header.Get("User-Agent"); ua != "" {
		return meta.String(ua)
	}

	if ua := m.UserAgent(ctx); ua != nil {
		return ua
	}

	return meta.String(userAgent)
}

func extractRequestID(ctx context.Context, req *http.Request) meta.Valuer {
	if id := req.Header.Get("Request-ID"); id != "" {
		return meta.String(id)
	}

	if id := m.RequestID(ctx); id != nil {
		return id
	}

	return meta.ToValuer(uuid.New())
}
