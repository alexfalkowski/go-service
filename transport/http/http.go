package http

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/net/http"
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

// Client is an alias for net/http.Client.
type Client = http.Client

// MaxBytesError is an alias for net/http.MaxBytesError.
type MaxBytesError = http.MaxBytesError

// Handler is an alias for net/http.Handler.
type Handler = http.Handler

// HandlerFunc is an alias for net/http.HandlerFunc.
type HandlerFunc = http.HandlerFunc

// Header is an alias for net/http.Header.
type Header = http.Header

// Request is an alias for net/http.Request.
type Request = http.Request

// Response is an alias for net/http.Response.
type Response = http.Response

// ResponseWriter is an alias for net/http.ResponseWriter.
type ResponseWriter = http.ResponseWriter

// RoundTripper is an alias for net/http.RoundTripper.
type RoundTripper = http.RoundTripper

// ServeMux is an alias for net/http.ServeMux.
type ServeMux = http.ServeMux

// DefaultTransport is an alias for http.DefaultTransport.
var DefaultTransport = http.DefaultTransport

// ErrUseLastResponse is an alias for http.ErrUseLastResponse.
var ErrUseLastResponse = http.ErrUseLastResponse

// ErrServerClosed is an alias for http.ErrServerClosed.
var ErrServerClosed = http.ErrServerClosed

// NoBody is an alias for http.NoBody.
var NoBody = http.NoBody

// NewRequestWithContext constructs a new outgoing HTTP request with ctx.
func NewRequestWithContext(ctx context.Context, method, url string, body io.Reader) (*Request, error) {
	return http.NewRequestWithContext(ctx, method, url, body)
}

// NewServeMux constructs a new HTTP request multiplexer.
func NewServeMux() *ServeMux {
	return http.NewServeMux()
}

// MaxBytesHandler wraps h so inbound request bodies are limited to n bytes.
func MaxBytesHandler(h Handler, n int64) Handler {
	return http.MaxBytesHandler(h, n)
}

// HandleFunc registers handler for pattern on mux and wraps it with OpenTelemetry instrumentation.
func HandleFunc(mux *ServeMux, pattern string, handler HandlerFunc) {
	http.HandleFunc(mux, pattern, handler)
}

// Handle registers handler for pattern on mux and wraps it with OpenTelemetry instrumentation.
func Handle(mux *ServeMux, pattern string, handler Handler) {
	http.Handle(mux, pattern, handler)
}

// StatusText returns the standard HTTP status text for the given status code.
func StatusText(code int) string {
	return http.StatusText(code)
}

// ParseServiceMethod derives a logical "service" and "method" name from an HTTP request.
func ParseServiceMethod(req *Request) (string, string) {
	return http.ParseServiceMethod(req)
}

// Pattern constructs a route pattern of the form "/<name><pattern>".
func Pattern(name env.Name, pattern string) string {
	return http.Pattern(name, pattern)
}
