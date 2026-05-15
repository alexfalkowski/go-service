package meta

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/net"
	"github.com/alexfalkowski/go-service/v2/net/header"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/net/http/strings"
	"github.com/alexfalkowski/go-service/v2/slices"
)

// NewHandler constructs server-side metadata middleware for HTTP requests.
//
// The returned handler extracts request metadata into the request context and sets standard response headers.
// It is designed to be used early in the server middleware chain so downstream middleware and handlers can
// rely on a populated context (for example, logging, auth, rate limiting, and tracing).
func NewHandler(userAgent env.UserAgent, version env.Version, generator id.Generator) *Handler {
	return &Handler{userAgent: userAgent, serviceVersion: version.String(), generator: generator}
}

// Handler extracts request metadata and stores it in the request context.
//
// Extracted metadata includes user agent, request id, client IP address (and its source kind), geolocation,
// and Authorization token value (when present and parseable).
type Handler struct {
	generator      id.Generator
	userAgent      env.UserAgent
	serviceVersion string
}

// ServeHTTP extracts metadata from req and stores it in the request context.
//
// Ignorable paths (health/metrics/etc.) bypass extraction.
//
// Response headers:
//   - "Service-Version" is set to the configured service version.
//   - "Request-Id" is set to the resolved request id.
//
// Context population:
//
// The handler populates the request context with:
//   - user agent (from context, request header, or default userAgent parameter)
//   - request id (from context, request header, or generated via generator)
//   - client IP address and IP address kind (derived from forwarded headers or RemoteAddr)
//   - geolocation (from "Geolocation" header)
//   - authorization value (derived from the "Authorization" header, when present)
//
// Error handling:
//
// If the Authorization header is present but cannot be parsed (unsupported scheme or invalid format),
// it writes an HTTP 400 error response and does not call next.
func (h *Handler) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if strings.IsIgnorable(req.URL.Path) {
		next(res, req)
		return
	}

	ctx := req.Context()
	userAgent := serverUserAgent(ctx, req, h.userAgent)

	requestID := serverRequestID(ctx, h.generator, req)
	serverSetResponseHeaders(res.Header(), h.serviceVersion, requestID.Value())

	kind, ip := serverIP(req)
	geolocation := serverGeolocation(req)

	auth, err := serverAuthorization(req)
	if err != nil {
		_ = status.WriteError(res, status.BadRequestError(err))
		return
	}
	ctx = meta.WithAttributes(ctx,
		meta.WithUserAgent(userAgent),
		meta.WithRequestID(requestID),
		meta.WithIPAddr(ip),
		meta.WithIPAddrKind(kind),
		meta.WithGeolocation(geolocation),
		meta.WithAuthorization(auth),
	)

	next(res, req.WithContext(ctx))
}

func serverSetResponseHeaders(header http.Header, serviceVersion, requestID string) {
	// Clip caps each header at one element so later appends allocate instead of
	// overwriting the neighboring value in this backing array.
	values := [...]string{serviceVersion, requestID}
	header["Service-Version"] = slices.Clip(values[0:1])
	header["Request-Id"] = slices.Clip(values[1:2])
}

func serverUserAgent(ctx context.Context, req *http.Request, userAgent env.UserAgent) meta.Value {
	if ua := meta.UserAgent(ctx); !ua.IsEmpty() {
		return ua
	}

	if ua := req.Header.Get("User-Agent"); !strings.IsEmpty(ua) {
		return meta.String(ua)
	}

	return meta.String(userAgent.String())
}

func serverRequestID(ctx context.Context, generator id.Generator, req *http.Request) meta.Value {
	if id := meta.RequestID(ctx); !id.IsEmpty() {
		return id
	}

	if id := req.Header.Get("Request-Id"); !strings.IsEmpty(id) {
		return meta.String(id)
	}

	return meta.String(generator.Generate())
}

func serverIP(req *http.Request) (meta.Value, meta.Value) {
	for _, h := range header.ForwardedIPs {
		if ip := req.Header.Get(h.HTTP); !strings.IsEmpty(ip) {
			ip, _, _ := strings.Cut(ip, ",")

			return meta.String(h.GRPC), meta.String(ip)
		}
	}

	remoteKind := meta.String("remote")
	addr := req.RemoteAddr

	return remoteKind, meta.String(net.Host(addr))
}

func serverAuthorization(req *http.Request) (meta.Value, error) {
	a := req.Header.Get("Authorization")
	if strings.IsEmpty(a) {
		return meta.Blank(), nil
	}

	value, err := header.ParseBearer(a)
	if err != nil {
		return meta.Blank(), err
	}

	return meta.Ignored(value), nil
}

func serverGeolocation(req *http.Request) meta.Value {
	return meta.String(req.Header.Get("Geolocation"))
}
