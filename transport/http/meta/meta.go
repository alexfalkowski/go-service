package meta

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/alexfalkowski/go-service/meta"
	snet "github.com/alexfalkowski/go-service/net"
	tmeta "github.com/alexfalkowski/go-service/transport/meta"
	"github.com/alexfalkowski/go-service/version"
	"github.com/google/uuid"
)

// Handler for meta.
type Handler struct {
	version version.Version
	http.Handler
}

// NewHandler for meta.
func NewHandler(version version.Version, handler http.Handler) *Handler {
	return &Handler{version: version, Handler: handler}
}

func (h *Handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	ctx = meta.WithVersion(ctx, string(h.version))
	ctx = tmeta.WithUserAgent(ctx, extractUserAgent(ctx, req))

	requestID := extractRequestID(ctx, req)
	if requestID == "" {
		requestID = uuid.New().String()
	}

	ctx = tmeta.WithRequestID(ctx, requestID)
	ctx = tmeta.WithRemoteAddress(ctx, extractRemoteAddress(ctx, req))

	resp.Header().Add("Version", string(h.version))

	h.Handler.ServeHTTP(resp, req.WithContext(ctx))
}

// NewRoundTripper for meta.
func NewRoundTripper(userAgent string, version version.Version, hrt http.RoundTripper) *RoundTripper {
	return &RoundTripper{userAgent: userAgent, version: version, RoundTripper: hrt}
}

// RoundTripper for meta.
type RoundTripper struct {
	userAgent string
	version   version.Version
	http.RoundTripper
}

func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	ctx = meta.WithVersion(ctx, string(r.version))

	userAgent := tmeta.UserAgent(ctx)
	if userAgent == "" {
		userAgent = r.userAgent
	}

	req.Header.Set("User-Agent", userAgent)
	ctx = tmeta.WithUserAgent(ctx, userAgent)

	requestID := tmeta.RequestID(ctx)
	if requestID == "" {
		requestID = uuid.New().String()
	}

	req.Header.Set("Request-ID", requestID)
	ctx = tmeta.WithRequestID(ctx, requestID)

	remoteAddress := tmeta.RemoteAddress(ctx)
	if remoteAddress == "" {
		remoteAddress = snet.OutboundIP(ctx)
	}

	req.Header.Set("X-Forwarded-For", remoteAddress)
	ctx = tmeta.WithRemoteAddress(ctx, remoteAddress)

	return r.RoundTripper.RoundTrip(req.WithContext(ctx))
}

func extractUserAgent(ctx context.Context, req *http.Request) string {
	if userAgent := req.Header.Get("User-Agent"); userAgent != "" {
		return userAgent
	}

	return tmeta.UserAgent(ctx)
}

func extractRequestID(ctx context.Context, req *http.Request) string {
	if requestID := req.Header.Get("Request-ID"); requestID != "" {
		return requestID
	}

	return tmeta.RequestID(ctx)
}

func extractRemoteAddress(ctx context.Context, req *http.Request) string {
	if forwardedFor := req.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		return strings.Split(forwardedFor, ",")[0]
	}

	if ip, _, err := net.SplitHostPort(req.RemoteAddr); err != nil && ip != "" {
		return ip
	}

	return tmeta.RemoteAddress(ctx)
}
