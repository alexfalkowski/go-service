package meta

import (
	"context"
	"net/http"

	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/id"
	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/net"
	"github.com/alexfalkowski/go-service/strings"
	"github.com/alexfalkowski/go-service/transport/header"
	tm "github.com/alexfalkowski/go-service/transport/meta"
	ts "github.com/alexfalkowski/go-service/transport/strings"
)

// NewHandler for meta.
func NewHandler(userAgent env.UserAgent, version env.Version, generator id.Generator) *Handler {
	return &Handler{userAgent: userAgent, version: version, generator: generator}
}

// Handler for meta.
type Handler struct {
	generator id.Generator
	userAgent env.UserAgent
	version   env.Version
}

func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if ts.IsObservable(req.URL.Path) {
		next(res, req)

		return
	}

	header := res.Header()
	header.Add("Service-Version", h.version.String())

	ctx := req.Context()
	ctx = tm.WithUserAgent(ctx, extractUserAgent(ctx, req, h.userAgent))

	requestID := extractRequestID(ctx, h.generator, req)

	header.Set("Request-Id", requestID.Value())
	ctx = tm.WithRequestID(ctx, requestID)

	kind, ip := extractIP(req)
	ctx = tm.WithIPAddr(ctx, ip)
	ctx = tm.WithIPAddrKind(ctx, kind)

	ctx = tm.WithGeolocation(ctx, extractGeolocation(req))
	ctx = tm.WithAuthorization(ctx, extractAuthorization(ctx, req))

	next(res, req.WithContext(ctx))
}

// NewRoundTripper for meta.
func NewRoundTripper(userAgent env.UserAgent, generator id.Generator, hrt http.RoundTripper) *RoundTripper {
	return &RoundTripper{userAgent: userAgent, generator: generator, RoundTripper: hrt}
}

// RoundTripper for meta.
type RoundTripper struct {
	generator id.Generator
	http.RoundTripper
	userAgent env.UserAgent
}

func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()

	userAgent := extractUserAgent(ctx, req, r.userAgent)

	req.Header.Set("User-Agent", userAgent.Value())
	ctx = tm.WithUserAgent(ctx, userAgent)

	requestID := extractRequestID(ctx, r.generator, req)

	req.Header.Set("Request-Id", requestID.Value())
	ctx = tm.WithRequestID(ctx, requestID)

	return r.RoundTripper.RoundTrip(req.WithContext(ctx))
}

func extractUserAgent(ctx context.Context, req *http.Request, userAgent env.UserAgent) meta.Value {
	if ua := tm.UserAgent(ctx); !ua.IsEmpty() {
		return ua
	}

	if ua := req.Header.Get("User-Agent"); !strings.IsEmpty(ua) {
		return meta.String(ua)
	}

	return meta.String(userAgent.String())
}

func extractRequestID(ctx context.Context, generator id.Generator, req *http.Request) meta.Value {
	if id := tm.RequestID(ctx); !id.IsEmpty() {
		return id
	}

	if id := req.Header.Get("Request-Id"); !strings.IsEmpty(id) {
		return meta.String(id)
	}

	return meta.String(generator.Generate())
}

func extractIP(req *http.Request) (meta.Value, meta.Value) {
	headers := []string{"X-Real-Ip", "CF-Connecting-Ip", "True-Client-Ip", "X-Forwarded-For"}
	for _, h := range headers {
		if ip := req.Header.Get(h); !strings.IsEmpty(ip) {
			ip, _, _ := strings.Cut(ip, ",")

			return meta.String(strings.ToLower(h)), meta.String(ip)
		}
	}

	remoteKind := meta.String("remote")
	addr := req.RemoteAddr

	return remoteKind, meta.String(net.Host(addr))
}

func extractAuthorization(ctx context.Context, req *http.Request) meta.Value {
	a := req.Header.Get("Authorization")
	if strings.IsEmpty(a) {
		return meta.Blank()
	}

	_, value, err := header.ParseAuthorization(a)
	if err != nil {
		meta.WithAttribute(ctx, "authError", meta.Error(err))

		return meta.Blank()
	}

	return meta.Ignored(value)
}

func extractGeolocation(req *http.Request) meta.Value {
	return meta.String(req.Header.Get("Geolocation"))
}
