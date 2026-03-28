package meta

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/net"
	"github.com/alexfalkowski/go-service/v2/net/header"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/net/http/strings"
)

const (
	requestKey  = context.Key("request")
	responseKey = context.Key("response")
	encoderKey  = context.Key("encoder")
)

// NoPrefix is an alias for meta.NoPrefix.
const NoPrefix = meta.NoPrefix

// Map is an alias for meta.Map.
type Map = meta.Map

// CamelStrings exports all stored meta attributes as a string map with lowerCamelCased keys.
//
// This is a thin wrapper around meta.CamelStrings. The prefix parameter is prepended to each exported key
// (if non-empty). Attributes whose rendered value is empty are skipped.
func CamelStrings(ctx context.Context, prefix string) Map {
	return meta.CamelStrings(ctx, prefix)
}

// Error converts err to a meta.Value using err.Error().
//
// This is a thin wrapper around meta.Error.
func Error(err error) meta.Value {
	return meta.Error(err)
}

// WithAttribute stores an arbitrary meta attribute on ctx.
//
// This is a thin wrapper around meta.WithAttribute.
func WithAttribute(ctx context.Context, key string, value meta.Value) context.Context {
	return meta.WithAttribute(ctx, key, value)
}

// WithRequest stores req in ctx and returns the derived context.
//
// This is commonly used by go-service HTTP content handlers/middleware to make the request available to
// downstream handlers via Request(ctx).
func WithRequest(ctx context.Context, req *http.Request) context.Context {
	return context.WithValue(ctx, requestKey, req)
}

// Request returns the stored *http.Request from ctx.
//
// Panics: Request expects WithRequest to have been called. It will panic if no request is stored in ctx
// or if the stored value is not a *http.Request.
func Request(ctx context.Context) *http.Request {
	return ctx.Value(requestKey).(*http.Request)
}

// WithResponse stores res in ctx and returns the derived context.
//
// This is commonly used by go-service HTTP content handlers/middleware to make the response writer available
// to downstream handlers via Response(ctx).
func WithResponse(ctx context.Context, res http.ResponseWriter) context.Context {
	return context.WithValue(ctx, responseKey, res)
}

// Response returns the stored http.ResponseWriter from ctx.
//
// Panics: Response expects WithResponse to have been called. It will panic if no response writer is stored
// in ctx or if the stored value is not an http.ResponseWriter.
func Response(ctx context.Context) http.ResponseWriter {
	return ctx.Value(responseKey).(http.ResponseWriter)
}

// WithEncoder stores enc in ctx and returns the derived context.
//
// This is commonly used by go-service HTTP content handlers/middleware to make the negotiated encoder
// (selected from Content-Type) available to downstream handlers via Encoder(ctx).
func WithEncoder(ctx context.Context, enc encoding.Encoder) context.Context {
	return context.WithValue(ctx, encoderKey, enc)
}

// Encoder returns the stored encoding.Encoder from ctx.
//
// Panics: Encoder expects WithEncoder to have been called. It will panic if no encoder is stored in ctx
// or if the stored value is not an encoding.Encoder.
func Encoder(ctx context.Context) encoding.Encoder {
	return ctx.Value(encoderKey).(encoding.Encoder)
}

// NewHandler constructs server-side metadata middleware for HTTP requests.
//
// The returned handler extracts request metadata into the request context and sets standard response headers.
// It is designed to be used early in the server middleware chain so downstream middleware and handlers can
// rely on a populated context (for example, logging, auth, rate limiting, and tracing).
func NewHandler(userAgent env.UserAgent, version env.Version, generator id.Generator) *Handler {
	return &Handler{userAgent: userAgent, version: version, generator: generator}
}

// Handler extracts request metadata and stores it in the request context.
//
// Extracted metadata includes user agent, request id, client IP address (and its source kind), geolocation,
// and Authorization token value (when present and parseable).
type Handler struct {
	generator id.Generator
	userAgent env.UserAgent
	version   env.Version
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

	header := res.Header()
	header.Add("Service-Version", h.version.String())

	ctx := req.Context()
	ctx = meta.WithUserAgent(ctx, extractUserAgent(ctx, req, h.userAgent))

	requestID := extractRequestID(ctx, h.generator, req)
	header.Set("Request-Id", requestID.Value())
	ctx = meta.WithRequestID(ctx, requestID)

	kind, ip := extractIP(req)
	ctx = meta.WithIPAddr(ctx, ip)
	ctx = meta.WithIPAddrKind(ctx, kind)

	ctx = meta.WithGeolocation(ctx, extractGeolocation(req))

	auth, err := extractAuthorization(req)
	if err != nil {
		status.WriteError(ctx, res, status.BadRequestError(err))
		return
	}
	ctx = meta.WithAuthorization(ctx, auth)

	next(res, req.WithContext(ctx))
}

// NewRoundTripper constructs client-side metadata middleware for HTTP requests.
//
// The returned RoundTripper injects standard request headers and synchronizes them back into the request
// context so downstream transport wrappers (for example logging/tracing) can read consistent values.
func NewRoundTripper(userAgent env.UserAgent, generator id.Generator, hrt http.RoundTripper) *RoundTripper {
	return &RoundTripper{RoundTripper: hrt, userAgent: userAgent, generator: generator}
}

// RoundTripper wraps an underlying `http.RoundTripper` and injects request metadata.
//
// This RoundTripper is intended to be applied as an outer wrapper so that other RoundTripper middleware
// (retry/breaker/logger/token, etc.) observes the finalized headers and context values.
type RoundTripper struct {
	http.RoundTripper
	generator id.Generator
	userAgent env.UserAgent
}

// RoundTrip injects request metadata into the outbound request.
//
// It sets the "User-Agent" and "Request-Id" headers, preferring values already present in the context or
// request headers, and stores the chosen values back into the request context.
//
// Precedence rules:
//   - If the context already contains a value (meta.UserAgent/meta.RequestID), that value is used.
//   - Else, if the request header already contains a value, that value is used.
//   - Else, a default is used (userAgent parameter or a generated request id).
func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()

	userAgent := extractUserAgent(ctx, req, r.userAgent)

	req.Header.Set("User-Agent", userAgent.Value())
	ctx = meta.WithUserAgent(ctx, userAgent)

	requestID := extractRequestID(ctx, r.generator, req)

	req.Header.Set("Request-Id", requestID.Value())
	ctx = meta.WithRequestID(ctx, requestID)

	return r.RoundTripper.RoundTrip(req.WithContext(ctx))
}

func extractUserAgent(ctx context.Context, req *http.Request, userAgent env.UserAgent) meta.Value {
	if ua := meta.UserAgent(ctx); !ua.IsEmpty() {
		return ua
	}

	if ua := req.Header.Get("User-Agent"); !strings.IsEmpty(ua) {
		return meta.String(ua)
	}

	return meta.String(userAgent.String())
}

func extractRequestID(ctx context.Context, generator id.Generator, req *http.Request) meta.Value {
	if id := meta.RequestID(ctx); !id.IsEmpty() {
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

func extractAuthorization(req *http.Request) (meta.Value, error) {
	a := req.Header.Get("Authorization")
	if strings.IsEmpty(a) {
		return meta.Blank(), nil
	}

	_, value, err := header.ParseAuthorization(a)
	if err != nil {
		return meta.Blank(), err
	}

	return meta.Ignored(value), nil
}

func extractGeolocation(req *http.Request) meta.Value {
	return meta.String(req.Header.Get("Geolocation"))
}
