package meta

import (
	"context"
	"net/http"

	"github.com/alexfalkowski/go-service/transport/meta"
	"github.com/google/uuid"
)

// Handler for meta.
type Handler struct{}

// NewHandler for meta.
func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) ServeHTTP(resp http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	ctx := req.Context()
	ctx = meta.WithUserAgent(ctx, extractUserAgent(ctx, req))

	requestID := extractRequestID(ctx, req)
	if requestID == "" {
		requestID = uuid.New().String()
	}

	ctx = meta.WithRequestID(ctx, requestID)

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

	userAgent := meta.UserAgent(ctx)
	if userAgent == "" {
		userAgent = r.userAgent
	}

	req.Header.Set("User-Agent", userAgent)
	ctx = meta.WithUserAgent(ctx, userAgent)

	requestID := meta.RequestID(ctx)
	if requestID == "" {
		requestID = uuid.New().String()
	}

	req.Header.Set("Request-ID", requestID)
	ctx = meta.WithRequestID(ctx, requestID)

	return r.RoundTripper.RoundTrip(req.WithContext(ctx))
}

func extractUserAgent(ctx context.Context, req *http.Request) string {
	if userAgent := req.Header.Get("User-Agent"); userAgent != "" {
		return userAgent
	}

	return meta.UserAgent(ctx)
}

func extractRequestID(ctx context.Context, req *http.Request) string {
	if requestID := req.Header.Get("Request-ID"); requestID != "" {
		return requestID
	}

	return meta.RequestID(ctx)
}
