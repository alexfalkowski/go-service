package meta

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/transport/header"
	m "github.com/alexfalkowski/go-service/transport/meta"
	ts "github.com/alexfalkowski/go-service/transport/strings"
	"github.com/google/uuid"
)

// Handler for meta.
type Handler struct {
	userAgent env.UserAgent
	version   env.Version
}

// NewHandler for meta.
func NewHandler(userAgent env.UserAgent, version env.Version) *Handler {
	return &Handler{userAgent: userAgent, version: version}
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

	requestID := extractRequestID(ctx, req)

	header.Set("Request-Id", requestID.Value())
	ctx = m.WithRequestID(ctx, requestID)

	kind, ip := extractIP(req)
	ctx = m.WithIPAddr(ctx, ip)
	ctx = m.WithIPAddrKind(ctx, kind)

	ctx = m.WithGeolocation(ctx, extractGeolocation(ctx, req))
	ctx = m.WithAuthorization(ctx, extractAuthorization(ctx, req))

	next(res, req.WithContext(ctx))
}

// NewRoundTripper for meta.
func NewRoundTripper(userAgent env.UserAgent, hrt http.RoundTripper) *RoundTripper {
	return &RoundTripper{userAgent: userAgent, RoundTripper: hrt}
}

// RoundTripper for meta.
type RoundTripper struct {
	http.RoundTripper
	userAgent env.UserAgent
}

func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()

	userAgent := extractUserAgent(ctx, req, r.userAgent)

	req.Header.Set("User-Agent", userAgent.Value())
	ctx = m.WithUserAgent(ctx, userAgent)

	requestID := extractRequestID(ctx, req)

	req.Header.Set("Request-Id", requestID.Value())
	ctx = m.WithRequestID(ctx, requestID)

	return r.RoundTripper.RoundTrip(req.WithContext(ctx))
}

func extractUserAgent(ctx context.Context, req *http.Request, userAgent env.UserAgent) meta.Valuer {
	if ua := m.UserAgent(ctx); ua != nil {
		return ua
	}

	if ua := req.Header.Get("User-Agent"); ua != "" {
		return meta.String(ua)
	}

	return meta.String(userAgent)
}

func extractRequestID(ctx context.Context, req *http.Request) meta.Valuer {
	if id := m.RequestID(ctx); id != nil {
		return id
	}

	if id := req.Header.Get("Request-Id"); id != "" {
		return meta.String(id)
	}

	return meta.ToString(uuid.New())
}

func extractIP(req *http.Request) (meta.Valuer, meta.Valuer) {
	headers := []string{"X-Real-Ip", "CF-Connecting-Ip", "True-Client-Ip", "X-Forwarded-For"}
	for _, h := range headers {
		if ip := req.Header.Get(h); ip != "" {
			ip, _, _ := strings.Cut(ip, ",")

			return meta.String(strings.ToLower(h)), meta.String(ip)
		}
	}

	remoteKind := meta.String("remote")
	addr := req.RemoteAddr

	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return remoteKind, meta.String(addr)
	}

	return remoteKind, meta.String(host)
}

func extractAuthorization(ctx context.Context, req *http.Request) meta.Valuer {
	a := req.Header.Get("Authorization")
	if a == "" {
		return meta.Blank()
	}

	_, t, err := header.ParseAuthorization(a)
	if err != nil {
		meta.WithAttribute(ctx, "authError", meta.Error(err))

		return meta.Blank()
	}

	return meta.Ignored(t)
}

func extractGeolocation(ctx context.Context, req *http.Request) meta.Valuer {
	if gl := m.Geolocation(ctx); gl != nil {
		return gl
	}

	return meta.String(req.Header.Get("Geolocation"))
}
