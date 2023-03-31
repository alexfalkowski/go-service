package meta

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/alexfalkowski/go-service/transport/meta"
	"github.com/google/uuid"
)

// Handler for meta.
type Handler struct {
	http.Handler
}

// NewHandler for meta.
func NewHandler(handler http.Handler) *Handler {
	return &Handler{Handler: handler}
}

func (h *Handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	ctx = meta.WithUserAgent(ctx, extractUserAgent(ctx, req))

	requestID := extractRequestID(ctx, req)
	if requestID == "" {
		requestID = uuid.New().String()
	}

	ctx = meta.WithRequestID(ctx, requestID)
	ctx = meta.WithRemoteAddress(ctx, extractRemoteAddress(ctx, req))

	h.Handler.ServeHTTP(resp, req.WithContext(ctx))
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

	ctx = meta.WithRemoteAddress(ctx, extractRemoteAddress(ctx, req))

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

func extractRemoteAddress(ctx context.Context, req *http.Request) string {
	if forwardedFor := req.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		return strings.Split(forwardedFor, ",")[0]
	}

	if host, _, err := net.SplitHostPort(req.RemoteAddr); err != nil && host != "" {
		return host
	}

	return meta.RemoteAddress(ctx)
}
