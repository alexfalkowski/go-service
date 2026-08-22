package http

import (
	"net/http"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/net/http/telemetry"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/urfave/negroni/v3"
)

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

// NewResponseController constructs a ResponseController for rw.
//
// This is a thin wrapper around [net/http.NewResponseController].
func NewResponseController(rw ResponseWriter) *ResponseController {
	return http.NewResponseController(rw)
}

// NewTelemetryHandler wraps handler with OpenTelemetry instrumentation when tracing or metrics are enabled.
func NewTelemetryHandler(handler Handler, operation string) Handler {
	if !metrics.IsEnabled() && !tracer.IsEnabled() {
		return handler
	}

	return telemetry.NewHandler(handler, operation)
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
