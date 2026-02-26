package telemetry

import (
	"context"
	"net/http"
	"net/http/httptrace"

	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// NewClientTrace returns a net/http/httptrace.ClientTrace that is integrated with
// OpenTelemetry.
//
// The returned trace can be attached to an outbound request context using
// httptrace.WithClientTrace. It observes HTTP client request lifecycle events
// (DNS start/done, connect start/done, TLS handshake, etc.) and records
// additional telemetry according to the provided options and the OpenTelemetry
// SDK configuration in your application.
//
// This function delegates to otelhttptrace.NewClientTrace. For exact semantics
// and supported otelhttptrace.ClientTraceOption values, consult the upstream
// otelhttptrace documentation for the version you vend.
func NewClientTrace(ctx context.Context, opts ...otelhttptrace.ClientTraceOption) *httptrace.ClientTrace {
	return otelhttptrace.NewClientTrace(ctx, opts...)
}

// WithClientTrace returns an otelhttp.Option that instructs an otelhttp
// instrumented transport to attach an httptrace.ClientTrace derived from a
// request context.
//
// This is typically used when constructing an instrumented transport (see
// NewTransport) so that each request can have an httptrace trace created from
// its context.
//
// This function delegates to otelhttp.WithClientTrace. For exact semantics and
// supported behavior, consult the upstream otelhttp documentation for the
// version you vend.
func WithClientTrace(f func(context.Context) *httptrace.ClientTrace) otelhttp.Option {
	return otelhttp.WithClientTrace(f)
}

// NewTransport returns an http.RoundTripper instrumented with OpenTelemetry.
//
// It wraps the provided base transport (for example http.DefaultTransport) and
// produces client-side spans and/or metrics for outbound HTTP requests according
// to the provided options and the OpenTelemetry SDK configuration in your
// application.
//
// The returned value is *otelhttp.Transport, which implements http.RoundTripper
// and can be installed on an http.Client.
//
// This function delegates to otelhttp.NewTransport. For exact semantics and
// supported otelhttp.Option values, consult the upstream otelhttp documentation
// for the version you vend.
func NewTransport(base http.RoundTripper, opts ...otelhttp.Option) *otelhttp.Transport {
	return otelhttp.NewTransport(base, opts...)
}

// NewHandler returns an http.Handler that is instrumented with OpenTelemetry.
//
// It wraps the provided handler and creates server-side spans and/or metrics for
// inbound HTTP requests according to the provided options and the OpenTelemetry
// SDK configuration in your application.
//
// The operation parameter is used by upstream otelhttp as the span name or a
// component of the span name (depending on the upstream version and options).
//
// This function delegates to otelhttp.NewHandler. For exact semantics and
// supported otelhttp.Option values, consult the upstream otelhttp documentation
// for the version you vend.
func NewHandler(handler http.Handler, operation string, opts ...otelhttp.Option) http.Handler {
	return otelhttp.NewHandler(handler, operation, opts...)
}
