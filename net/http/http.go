package http

import (
	"net/http"
	"net/http/httptrace"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/config/options"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/net/http/telemetry"
	"github.com/alexfalkowski/go-service/v2/net/url"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/urfave/negroni/v3"
)

// DefaultMaxHeaderBytes is an alias of [http.DefaultMaxHeaderBytes].
const DefaultMaxHeaderBytes = http.DefaultMaxHeaderBytes

// TimeFormat is the HTTP time format layout.
//
// It is an alias of [http.TimeFormat].
const TimeFormat = http.TimeFormat

// Client is an alias for [net/http.Client].
//
// It is provided so go-service code can depend on a consistent import path while preserving
// standard library semantics.
type Client = http.Client

// MaxBytesError is an alias for [net/http.MaxBytesError].
//
// It is returned when MaxBytesReader or MaxBytesHandler observes an inbound request body exceeding
// the configured byte limit.
type MaxBytesError = http.MaxBytesError

// Handler is an alias for [net/http.Handler].
//
// It is provided so go-service code can depend on a consistent import path while preserving
// standard library semantics.
type Handler = http.Handler

// ChainedHandler is an alias for negroni.Handler.
type ChainedHandler = negroni.Handler

// ChainedHandlers is an alias for negroni.Negroni.
type ChainedHandlers = negroni.Negroni

// HandlerFunc is an alias for [net/http.HandlerFunc].
//
// It is provided so go-service code can depend on a consistent import path while preserving
// standard library semantics.
type HandlerFunc = http.HandlerFunc

// Header is an alias for [net/http.Header].
//
// It is provided so go-service code can depend on a consistent import path while preserving
// standard library semantics.
type Header = http.Header

// Request is an alias for [net/http.Request].
//
// It is provided so go-service code can depend on a consistent import path while preserving
// standard library semantics.
type Request = http.Request

// Response is an alias for [net/http.Response].
//
// It is provided so go-service code can depend on a consistent import path while preserving
// standard library semantics.
type Response = http.Response

// ServeMux is an alias for [net/http.ServeMux].
//
// It is provided so go-service code can depend on a consistent import path while preserving
// standard library semantics.
type ServeMux = http.ServeMux

// Server is an alias for [net/http.Server].
//
// It is provided so go-service code can depend on a consistent import path while preserving
// standard library semantics.
type Server = http.Server

// ResponseWriter is an alias for [net/http.ResponseWriter].
//
// It is provided so go-service code can depend on a consistent import path while preserving
// standard library semantics.
type ResponseWriter = http.ResponseWriter

// RoundTripper is an alias for [net/http.RoundTripper].
//
// It is provided so go-service code can depend on a consistent import path while preserving
// standard library semantics.
type RoundTripper = http.RoundTripper

// ClosingRoundTripper adapts a function to RoundTripper while making request-body ownership explicit.
//
// The function's third return value reports whether ClosingRoundTripper should close req.Body before returning.
// Return true when the function rejects the request locally without delegating to another RoundTripper.
// Return false after delegating because the delegated RoundTripper owns the request body.
type ClosingRoundTripper func(req *Request) (*Response, error, bool)

// RoundTrip calls s and closes req.Body when s asks it to.
func (s ClosingRoundTripper) RoundTrip(req *Request) (*Response, error) {
	res, err, closeBody := s(req)
	if closeBody && req != nil && req.Body != nil && req.Body != NoBody {
		_ = req.Body.Close()
	}

	return res, err
}

// DefaultTransport is an alias for [http.DefaultTransport].
var DefaultTransport = http.DefaultTransport

// ErrUseLastResponse is an alias for [http.ErrUseLastResponse].
var ErrUseLastResponse = http.ErrUseLastResponse

// ErrServerClosed is an alias for [http.ErrServerClosed].
var ErrServerClosed = http.ErrServerClosed

// NoBody is an alias for [http.NoBody].
var NoBody = http.NoBody

// NewClient constructs an HTTP client with a request timeout.
//
// When tracing or metrics are enabled, the returned client wraps the provided RoundTripper with a telemetry
// transport. When tracing is enabled, it also installs an httptrace-based client trace derived from the
// request context.
//
// The provided timeout is assigned to [http.Client.Timeout] (total time limit for requests, including
// connection time, redirects, and reading the response body).
func NewClient(rt http.RoundTripper, timeout time.Duration) *http.Client {
	var transport http.RoundTripper

	if metrics.IsEnabled() || tracer.IsEnabled() {
		options := []telemetry.Option{}
		if tracer.IsEnabled() {
			options = append(options, telemetry.WithClientTrace(func(ctx context.Context) *httptrace.ClientTrace {
				return telemetry.NewClientTrace(ctx)
			}))
		}

		transport = telemetry.NewTransport(rt, options...)
	} else {
		transport = rt
	}

	return &http.Client{
		Transport: transport,
		Timeout:   timeout.Duration(),
	}
}

// IgnoreRedirect returns redirect responses without following them.
func IgnoreRedirect(_ *Request, _ []*Request) error {
	return ErrUseLastResponse
}

// SameOrigin reports whether prev and next use the same URL origin.
func SameOrigin(prev, next *url.URL) bool {
	if prev == nil || next == nil {
		return false
	}

	return strings.ToLower(prev.Scheme) == strings.ToLower(next.Scheme) &&
		strings.ToLower(prev.Hostname()) == strings.ToLower(next.Hostname()) &&
		originPort(prev) == originPort(next)
}

func originPort(u *url.URL) string {
	if port := u.Port(); !strings.IsEmpty(port) {
		return port
	}

	switch strings.ToLower(u.Scheme) {
	case "http":
		return "80"
	case "https":
		return "443"
	default:
		return strings.Empty
	}
}

// IsCrossOriginRedirect reports whether req is a redirected request whose previous request used a different origin.
func IsCrossOriginRedirect(req *Request) bool {
	if req == nil || req.Response == nil || req.Response.Request == nil {
		return false
	}

	return !SameOrigin(req.Response.Request.URL, req.URL)
}

// SameOriginRedirect allows redirects only when the next request stays on the same scheme and host.
//
// It is intended for clients that add credentials or signatures in RoundTripper middleware, where Go's
// default cross-origin sensitive-header stripping is not enough because middleware may mint fresh credentials
// for each redirected request.
func SameOriginRedirect(req *Request, via []*Request) error {
	if len(via) == 0 {
		return nil
	}

	prev := via[len(via)-1].URL
	next := req.URL
	if SameOrigin(prev, next) {
		return nil
	}

	return ErrUseLastResponse
}

// NewRequestWithContext constructs a new outgoing HTTP request with ctx.
//
// This is a thin wrapper around [net/http.NewRequestWithContext]. The returned request is canceled
// when ctx is canceled.
func NewRequestWithContext(ctx context.Context, method, url string, body io.Reader) (*Request, error) {
	return http.NewRequestWithContext(ctx, method, url, body)
}

// NewServeMux constructs a new HTTP request multiplexer.
//
// This is a thin wrapper around [net/http.NewServeMux].
func NewServeMux() *ServeMux {
	return http.NewServeMux()
}

// NewChainedHandlers constructs an HTTP middleware chain.
func NewChainedHandlers(handlers ...ChainedHandler) *ChainedHandlers {
	return negroni.New(handlers...)
}

// MaxBytesHandler wraps h so inbound request bodies are limited to n bytes.
//
// This is a thin wrapper around [net/http.MaxBytesHandler].
func MaxBytesHandler(h Handler, n int64) Handler {
	return http.MaxBytesHandler(h, n)
}

// NewTelemetryHandler wraps handler with OpenTelemetry instrumentation when tracing or metrics are enabled.
func NewTelemetryHandler(handler Handler, operation string) Handler {
	if !metrics.IsEnabled() && !tracer.IsEnabled() {
		return handler
	}

	return telemetry.NewHandler(handler, operation)
}

// HandleFunc registers handler for pattern on mux.
func HandleFunc(mux *ServeMux, pattern string, handler http.HandlerFunc) {
	Handle(mux, pattern, handler)
}

// Handle registers handler for pattern on mux.
func Handle(mux *ServeMux, pattern string, handler http.Handler) {
	mux.Handle(pattern, handler)
}

// StatusText returns the standard HTTP status text for the given status code.
//
// This is a thin wrapper around [net/http.StatusText].
func StatusText(code int) string {
	return http.StatusText(code)
}

// ParseTime parses an HTTP time value.
//
// This is a thin wrapper around [http.ParseTime] and does not change semantics.
func ParseTime(value string) (time.Time, error) {
	return http.ParseTime(value)
}

// NewServer constructs an HTTP server with common timeout defaults and supported protocol settings.
//
// Timeouts are derived from options first (if present) and fall back to the provided timeout value:
//   - read_timeout
//   - write_timeout
//   - idle_timeout
//   - read_header_timeout
//
// Additional low-level server tuning may be provided through options using:
//   - max_header_bytes
//
// Protocols are configured via Protocols().
//
// Note: [options.NonNegativeDuration] uses MustParseDuration under the hood; invalid or negative option
// values will panic at server construction time.
func NewServer(options options.Map, timeout time.Duration, handler Handler) *Server {
	return &http.Server{
		Handler:           handler,
		ReadTimeout:       options.NonNegativeDuration("read_timeout", timeout).Duration(),
		WriteTimeout:      options.NonNegativeDuration("write_timeout", timeout).Duration(),
		IdleTimeout:       options.NonNegativeDuration("idle_timeout", timeout).Duration(),
		ReadHeaderTimeout: options.NonNegativeDuration("read_header_timeout", timeout).Duration(),
		MaxHeaderBytes:    options.IntSize("max_header_bytes", bytes.Size(DefaultMaxHeaderBytes)),
		Protocols:         Protocols(),
	}
}

// ParseServiceMethod derives a logical "service" and "method" name from an HTTP request.
//
// This helper is intended for consistent telemetry naming. It attempts to derive names from the request
// path when it follows the conventional go-service route shape:
//
//	/<service>/<method>
//
// If the request path matches that shape, ParseServiceMethod returns the extracted service/method pair.
//
// Otherwise it falls back to:
//   - method: lower-cased HTTP method (e.g. "get", "post")
//   - service: a best-effort name derived from the path:
//   - "root" when the path is empty or "/"
//   - otherwise the path without the leading "/" (e.g. "/health" -> "health")
func ParseServiceMethod(req *http.Request) (string, string) {
	path := req.URL.Path
	if service, method, ok := url.SplitPath(path); ok {
		return service, method
	}

	method := strings.ToLower(req.Method)

	if strings.IsEmpty(path) {
		return "root", method
	}

	path = path[1:]
	if strings.IsEmpty(path) {
		return "root", method
	}

	return path, method
}

// Pattern constructs a route pattern of the form "/<name><pattern>".
//
// This helper is used to namespace routes by service name so different services can share a router/mux
// without colliding, and so route names are consistent across telemetry, server registration, and tests.
//
// Example:
//
//	Pattern(name, "/debug/pprof/") // -> "/my-service/debug/pprof/"
func Pattern(name env.Name, pattern string) string {
	return strings.Concat("/", name.String(), pattern)
}
