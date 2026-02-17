package http

import (
	"net/http"
	"net/http/httptrace"

	"github.com/alexfalkowski/go-service/v2/config/options"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/transport/strings"
	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
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
	// Client is an alias for http.Client.
	Client = http.Client

	// Handler is an alias for http.Handler.
	Handler = http.Handler

	// HandlerFunc is an alias for http.HandlerFunc.
	HandlerFunc = http.HandlerFunc

	// Header is an alias for http.Header.
	Header = http.Header

	// Request is an alias for http.Request.
	Request = http.Request

	// Response is an alias for http.Response.
	Response = http.Response

	// ServeMux is an alias for http.ServeMux.
	ServeMux = http.ServeMux

	// Server is an alias for http.Server.
	Server = http.Server

	// ResponseWriter is an alias for http.ResponseWriter.
	ResponseWriter = http.ResponseWriter

	// RoundTripper is an alias for http.RoundTripper.
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

// NewClient returns an http.Client that instruments requests with OpenTelemetry and applies timeout.
func NewClient(rt http.RoundTripper, timeout time.Duration) *http.Client {
	return &http.Client{
		Transport: otelhttp.NewTransport(
			rt,
			otelhttp.WithClientTrace(func(ctx context.Context) *httptrace.ClientTrace {
				return otelhttptrace.NewClientTrace(ctx)
			}),
		),
		Timeout: timeout,
	}
}

// NewRequestWithContext is an alias for http.NewRequestWithContext.
func NewRequestWithContext(ctx context.Context, method, url string, body io.Reader) (*Request, error) {
	return http.NewRequestWithContext(ctx, method, url, body)
}

// NewServeMux is an alias for http.NewServeMux.
func NewServeMux() *ServeMux {
	return http.NewServeMux()
}

// HandleFunc registers handler for pattern and wraps it with OpenTelemetry instrumentation.
func HandleFunc(mux *ServeMux, pattern string, handler http.HandlerFunc) {
	Handle(mux, pattern, handler)
}

// Handle registers handler for pattern and wraps it with OpenTelemetry instrumentation.
func Handle(mux *ServeMux, pattern string, handler http.Handler) {
	mux.Handle(pattern, otelhttp.NewHandler(handler, pattern))
}

// StatusText is an alias for http.StatusText.
func StatusText(code int) string {
	return http.StatusText(code)
}

// NewServer builds an http.Server configured with common timeouts and supported protocols.
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

// ParseServiceMethod derives a service name and method from the request path and method.
//
// It uses transport/strings to split /service/method paths and falls back to the HTTP method
// and request path when no service/method pair is present.
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

// Pattern constructs a route pattern of the form /<name><pattern>.
func Pattern(name env.Name, pattern string) string {
	return strings.Concat("/", name.String(), pattern)
}
