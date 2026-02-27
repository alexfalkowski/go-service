package http

import (
	"net/http"
	"net/http/httptrace"

	"github.com/alexfalkowski/go-service/v2/config/options"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/net/http/telemetry"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/transport/strings"
)

const (
	// MethodDelete is an alias of http.MethodDelete.
	MethodDelete = http.MethodDelete

	// MethodGet is an alias of http.MethodGet.
	MethodGet = http.MethodGet

	// MethodPatch is an alias of http.MethodPatch.
	MethodPatch = http.MethodPatch

	// MethodPost is an alias of http.MethodPost.
	MethodPost = http.MethodPost

	// MethodPut is an alias of http.MethodPut.
	MethodPut = http.MethodPut

	// StatusBadRequest is an alias of http.StatusBadRequest.
	StatusBadRequest = http.StatusBadRequest

	// StatusConflict is an alias of http.StatusConflict.
	StatusConflict = http.StatusConflict

	// StatusForbidden is an alias of http.StatusForbidden.
	StatusForbidden = http.StatusForbidden

	// StatusGatewayTimeout is an alias of http.StatusGatewayTimeout.
	StatusGatewayTimeout = http.StatusGatewayTimeout

	// StatusOK is an alias of http.StatusOK.
	StatusOK = http.StatusOK

	// StatusInternalServerError is an alias of http.StatusInternalServerError.
	StatusInternalServerError = http.StatusInternalServerError

	// StatusNotFound is an alias of http.StatusNotFound.
	StatusNotFound = http.StatusNotFound

	// StatusNotImplemented is an alias of http.StatusNotImplemented.
	StatusNotImplemented = http.StatusNotImplemented

	// StatusServiceUnavailable is an alias of http.StatusServiceUnavailable.
	StatusServiceUnavailable = http.StatusServiceUnavailable

	// StatusTooManyRequests is an alias of http.StatusTooManyRequests.
	StatusTooManyRequests = http.StatusTooManyRequests

	// StatusUnauthorized is an alias of http.StatusUnauthorized.
	StatusUnauthorized = http.StatusUnauthorized
)

type (
	// Client is an alias for net/http.Client.
	//
	// It is provided so go-service code can depend on a consistent import path while preserving
	// standard library semantics.
	Client = http.Client
	// Handler is an alias for net/http.Handler.
	//
	// It is provided so go-service code can depend on a consistent import path while preserving
	// standard library semantics.
	Handler = http.Handler

	// HandlerFunc is an alias for net/http.HandlerFunc.
	//
	// It is provided so go-service code can depend on a consistent import path while preserving
	// standard library semantics.
	HandlerFunc = http.HandlerFunc

	// Header is an alias for net/http.Header.
	//
	// It is provided so go-service code can depend on a consistent import path while preserving
	// standard library semantics.
	Header = http.Header

	// Request is an alias for net/http.Request.
	//
	// It is provided so go-service code can depend on a consistent import path while preserving
	// standard library semantics.
	Request = http.Request

	// Response is an alias for net/http.Response.
	//
	// It is provided so go-service code can depend on a consistent import path while preserving
	// standard library semantics.
	Response = http.Response

	// ServeMux is an alias for net/http.ServeMux.
	//
	// It is provided so go-service code can depend on a consistent import path while preserving
	// standard library semantics.
	ServeMux = http.ServeMux

	// Server is an alias for net/http.Server.
	//
	// It is provided so go-service code can depend on a consistent import path while preserving
	// standard library semantics.
	Server = http.Server

	// ResponseWriter is an alias for net/http.ResponseWriter.
	//
	// It is provided so go-service code can depend on a consistent import path while preserving
	// standard library semantics.
	ResponseWriter = http.ResponseWriter

	// RoundTripper is an alias for net/http.RoundTripper.
	//
	// It is provided so go-service code can depend on a consistent import path while preserving
	// standard library semantics.
	RoundTripper = http.RoundTripper
)

var (
	// DefaultTransport is an alias for http.DefaultTransport.
	DefaultTransport = http.DefaultTransport

	// ErrUseLastResponse is an alias for http.ErrUseLastResponse.
	ErrUseLastResponse = http.ErrUseLastResponse

	// ErrServerClosed is an alias for http.ErrServerClosed.
	ErrServerClosed = http.ErrServerClosed

	// NoBody is an alias for http.NoBody.
	NoBody = http.NoBody
)

// NewClient constructs an HTTP client with OpenTelemetry instrumentation and a request timeout.
//
// The returned client wraps the provided RoundTripper with a telemetry transport and installs an
// httptrace-based client trace derived from the request context. This enables client-side spans and
// timing events to be captured by the configured OpenTelemetry instrumentation.
//
// The provided timeout is assigned to http.Client.Timeout (total time limit for requests, including
// connection time, redirects, and reading the response body).
func NewClient(rt http.RoundTripper, timeout time.Duration) *http.Client {
	return &http.Client{
		Transport: telemetry.NewTransport(
			rt,
			telemetry.WithClientTrace(func(ctx context.Context) *httptrace.ClientTrace {
				return telemetry.NewClientTrace(ctx)
			}),
		),
		Timeout: timeout,
	}
}

// NewRequestWithContext constructs a new outgoing HTTP request with ctx.
//
// This is a thin wrapper around net/http.NewRequestWithContext. The returned request is canceled
// when ctx is canceled.
func NewRequestWithContext(ctx context.Context, method, url string, body io.Reader) (*Request, error) {
	return http.NewRequestWithContext(ctx, method, url, body)
}

// NewServeMux constructs a new HTTP request multiplexer.
//
// This is a thin wrapper around net/http.NewServeMux.
func NewServeMux() *ServeMux {
	return http.NewServeMux()
}

// HandleFunc registers handler for pattern on mux and wraps it with OpenTelemetry instrumentation.
//
// This helper ensures that handlers registered via this package are consistently instrumented by
// wrapping them with telemetry.NewHandler before registration.
func HandleFunc(mux *ServeMux, pattern string, handler http.HandlerFunc) {
	Handle(mux, pattern, handler)
}

// Handle registers handler for pattern on mux and wraps it with OpenTelemetry instrumentation.
//
// The handler is wrapped with telemetry.NewHandler using the provided pattern as the handler name.
// This is useful for consistent HTTP server span naming and handler metrics attribution.
func Handle(mux *ServeMux, pattern string, handler http.Handler) {
	mux.Handle(pattern, telemetry.NewHandler(handler, pattern))
}

// StatusText returns the standard HTTP status text for the given status code.
//
// This is a thin wrapper around net/http.StatusText.
func StatusText(code int) string {
	return http.StatusText(code)
}

// NewServer constructs an HTTP server with common timeout defaults and supported protocol settings.
//
// Timeouts are derived from options first (if present) and fall back to the provided timeout value:
//   - read_timeout
//   - write_timeout
//   - idle_timeout
//   - read_header_timeout
//
// Protocols are configured via Protocols().
//
// Note: options.Duration uses MustParseDuration under the hood; invalid option values will panic at
// server construction time.
func NewServer(options options.Map, timeout time.Duration, handler Handler) *Server {
	return &http.Server{
		Handler:           handler,
		ReadTimeout:       options.Duration("read_timeout", timeout),
		WriteTimeout:      options.Duration("write_timeout", timeout),
		IdleTimeout:       options.Duration("idle_timeout", timeout),
		ReadHeaderTimeout: options.Duration("read_header_timeout", timeout),
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
// If the request path matches that shape (as determined by transport/strings.SplitServiceMethod),
// ParseServiceMethod returns the extracted service/method pair.
//
// Otherwise it falls back to:
//   - method: lower-cased HTTP method (e.g. "get", "post")
//   - service: a best-effort name derived from the path:
//   - "root" when the path is empty or "/"
//   - otherwise the path without the leading "/" (e.g. "/health" -> "health")
func ParseServiceMethod(req *http.Request) (string, string) {
	path := req.URL.Path
	if service, method, ok := strings.SplitServiceMethod(path); ok {
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
