package http

import (
	"fmt"
	"math"
	"net/http"
	"net/http/httptrace"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/config/options"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/net/grpc/strings"
	"github.com/alexfalkowski/go-service/v2/net/http/telemetry"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/alexfalkowski/go-service/v2/time"
)

// MethodDelete is an alias of http.MethodDelete.
const MethodDelete = http.MethodDelete

// MethodGet is an alias of http.MethodGet.
const MethodGet = http.MethodGet

// MethodPatch is an alias of http.MethodPatch.
const MethodPatch = http.MethodPatch

// MethodPost is an alias of http.MethodPost.
const MethodPost = http.MethodPost

// MethodPut is an alias of http.MethodPut.
const MethodPut = http.MethodPut

// StatusBadRequest is an alias of http.StatusBadRequest.
const StatusBadRequest = http.StatusBadRequest

// StatusConflict is an alias of http.StatusConflict.
const StatusConflict = http.StatusConflict

// StatusForbidden is an alias of http.StatusForbidden.
const StatusForbidden = http.StatusForbidden

// StatusGatewayTimeout is an alias of http.StatusGatewayTimeout.
const StatusGatewayTimeout = http.StatusGatewayTimeout

// StatusOK is an alias of http.StatusOK.
const StatusOK = http.StatusOK

// StatusRequestEntityTooLarge is an alias of http.StatusRequestEntityTooLarge.
const StatusRequestEntityTooLarge = http.StatusRequestEntityTooLarge

// StatusInternalServerError is an alias of http.StatusInternalServerError.
const StatusInternalServerError = http.StatusInternalServerError

// StatusMethodNotAllowed is an alias of http.StatusMethodNotAllowed.
const StatusMethodNotAllowed = http.StatusMethodNotAllowed

// StatusNotFound is an alias of http.StatusNotFound.
const StatusNotFound = http.StatusNotFound

// StatusNotImplemented is an alias of http.StatusNotImplemented.
const StatusNotImplemented = http.StatusNotImplemented

// StatusServiceUnavailable is an alias of http.StatusServiceUnavailable.
const StatusServiceUnavailable = http.StatusServiceUnavailable

// StatusTooManyRequests is an alias of http.StatusTooManyRequests.
const StatusTooManyRequests = http.StatusTooManyRequests

// StatusUnauthorized is an alias of http.StatusUnauthorized.
const StatusUnauthorized = http.StatusUnauthorized

// DefaultMaxHeaderBytes is an alias of http.DefaultMaxHeaderBytes.
const DefaultMaxHeaderBytes = http.DefaultMaxHeaderBytes

type (
	// Client is an alias for net/http.Client.
	//
	// It is provided so go-service code can depend on a consistent import path while preserving
	// standard library semantics.
	Client = http.Client

	// MaxBytesError is an alias for net/http.MaxBytesError.
	//
	// It is returned when MaxBytesReader or MaxBytesHandler observes an inbound request body exceeding
	// the configured byte limit.
	MaxBytesError = http.MaxBytesError
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

// NewClient constructs an HTTP client with a request timeout.
//
// When tracing or metrics are enabled, the returned client wraps the provided RoundTripper with a telemetry
// transport. When tracing is enabled, it also installs an httptrace-based client trace derived from the
// request context.
//
// The provided timeout is assigned to http.Client.Timeout (total time limit for requests, including
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

// MaxBytesHandler wraps h so inbound request bodies are limited to n bytes.
//
// This is a thin wrapper around net/http.MaxBytesHandler.
func MaxBytesHandler(h Handler, n int64) Handler {
	return http.MaxBytesHandler(h, n)
}

// HandleFunc registers handler for pattern on mux.
//
// When tracing or metrics are enabled, handler is wrapped with OpenTelemetry instrumentation before
// registration.
func HandleFunc(mux *ServeMux, pattern string, handler http.HandlerFunc) {
	Handle(mux, pattern, handler)
}

// Handle registers handler for pattern on mux.
//
// When tracing or metrics are enabled, handler is wrapped with telemetry.NewHandler using the provided
// pattern as the handler name.
func Handle(mux *ServeMux, pattern string, handler http.Handler) {
	if !metrics.IsEnabled() && !tracer.IsEnabled() {
		mux.Handle(pattern, handler)
		return
	}

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
// Additional low-level server tuning may be provided through options using:
//   - max_header_bytes
//
// Protocols are configured via Protocols().
//
// Note: options.Duration uses MustParseDuration under the hood; invalid option values will panic at
// server construction time.
func NewServer(options options.Map, timeout time.Duration, handler Handler) *Server {
	return &http.Server{
		Handler:           handler,
		ReadTimeout:       options.Duration("read_timeout", timeout).Duration(),
		WriteTimeout:      options.Duration("write_timeout", timeout).Duration(),
		IdleTimeout:       options.Duration("idle_timeout", timeout).Duration(),
		ReadHeaderTimeout: options.Duration("read_header_timeout", timeout).Duration(),
		MaxHeaderBytes:    mustIntSize(options, "max_header_bytes", bytes.Size(DefaultMaxHeaderBytes)),
		Protocols:         Protocols(),
	}
}

func mustIntSize(options options.Map, key string, fallback bytes.Size) int {
	size := options.Size(key, fallback)
	if size.Bytes() > math.MaxInt {
		runtime.Must(fmt.Errorf("http: %s exceeds max int: %s", key, size))
	}

	return int(size.Bytes())
}

// ParseServiceMethod derives a logical "service" and "method" name from an HTTP request.
//
// This helper is intended for consistent telemetry naming. It attempts to derive names from the request
// path when it follows the conventional go-service route shape:
//
//	/<service>/<method>
//
// If the request path matches that shape (as determined by net/grpc/strings.SplitServiceMethod),
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
