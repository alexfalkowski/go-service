package meta

import (
	"context"
	"net/http"
	"strings"

	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/id"
	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/net"
	"github.com/alexfalkowski/go-service/transport/header"
	m "github.com/alexfalkowski/go-service/transport/meta"
	ts "github.com/alexfalkowski/go-service/transport/strings"
)

// Handler for meta.
type Handler struct {
	gen       id.Generator
	userAgent env.UserAgent
	version   env.Version
}

// NewHandler for meta.
func NewHandler(userAgent env.UserAgent, version env.Version, gen id.Generator) *Handler {
	return &Handler{userAgent: userAgent, version: version, gen: gen}
}

func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if ts.IsObservable(req.URL.Path) {
		next(res, req)

		return
	}

	header := res.Header()
	header.Add("Service-Version", h.version.String())

	ctx := req.Context()
	ctx = m.WithUserAgent(ctx, extractUserAgent(ctx, req, h.userAgent))

	requestID := extractRequestID(ctx, h.gen, req)

	header.Set("Request-Id", requestID.Value())
	ctx = m.WithRequestID(ctx, requestID)

	kind, ip := extractIP(req)
	ctx = m.WithIPAddr(ctx, ip)
	ctx = m.WithIPAddrKind(ctx, kind)

	ctx = m.WithGeolocation(ctx, extractGeolocation(req))
	ctx = m.WithAuthorization(ctx, extractAuthorization(ctx, req))

	next(res, req.WithContext(ctx))
}

// NewRoundTripper for meta.
func NewRoundTripper(userAgent env.UserAgent, gen id.Generator, hrt http.RoundTripper) *RoundTripper {
	return &RoundTripper{userAgent: userAgent, gen: gen, RoundTripper: hrt}
}

// RoundTripper for meta.
type RoundTripper struct {
	gen id.Generator
	http.RoundTripper
	userAgent env.UserAgent
}

func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()

	userAgent := extractUserAgent(ctx, req, r.userAgent)

	req.Header.Set("User-Agent", userAgent.Value())
	ctx = m.WithUserAgent(ctx, userAgent)

	requestID := extractRequestID(ctx, r.gen, req)

	req.Header.Set("Request-Id", requestID.Value())
	ctx = m.WithRequestID(ctx, requestID)

	return r.RoundTripper.RoundTrip(req.WithContext(ctx))
}

func extractUserAgent(ctx context.Context, req *http.Request, userAgent env.UserAgent) *meta.Value {
	if ua := m.UserAgent(ctx); ua != nil {
		return ua
	}

	if ua := req.Header.Get("User-Agent"); ua != "" {
		return meta.String(ua)
	}

	return meta.String(string(userAgent))
}

func extractRequestID(ctx context.Context, gen id.Generator, req *http.Request) *meta.Value {
	if id := m.RequestID(ctx); id != nil {
		return id
	}

	if id := req.Header.Get("Request-Id"); id != "" {
		return meta.String(id)
	}

	return meta.String(gen.Generate())
}

func extractIP(req *http.Request) (*meta.Value, *meta.Value) {
	headers := []string{"X-Real-Ip", "CF-Connecting-Ip", "True-Client-Ip", "X-Forwarded-For"}
	for _, h := range headers {
		if ip := req.Header.Get(h); ip != "" {
			ip, _, _ := strings.Cut(ip, ",")

			return meta.String(strings.ToLower(h)), meta.String(ip)
		}
	}

	remoteKind := meta.String("remote")
	addr := req.RemoteAddr

	return remoteKind, meta.String(net.Host(addr))
}

func extractAuthorization(ctx context.Context, req *http.Request) *meta.Value {
	a := req.Header.Get("Authorization")
	if a == "" {
		return meta.Blank()
	}

	_, value, err := header.ParseAuthorization(a)
	if err != nil {
		meta.WithAttribute(ctx, "authError", meta.Error(err))

		return meta.Blank()
	}

	return meta.Ignored(value)
}

func extractGeolocation(req *http.Request) *meta.Value {
	return meta.String(req.Header.Get("Geolocation"))
}
