// Package telemetry provides minimal, stable helpers for wiring OpenTelemetry
// instrumentation into net/http clients and servers.
//
// This package is intentionally thin and delegates to the upstream OpenTelemetry
// HTTP instrumentation packages:
//
//   - go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp
//   - go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace
//
// The goal is to standardize how services attach HTTP telemetry without
// re-exporting the entire upstream API surface.
//
// # Server-side instrumentation
//
// Use NewHandler to wrap an existing http.Handler. The wrapper creates spans for
// incoming requests and records attributes and status information according to
// upstream otelhttp behavior and the provided options.
//
// Example:
//
//	mux := http.NewServeMux()
//	mux.Handle("/health", healthHandler)
//
//	h := telemetry.NewHandler(mux, "http.server")
//
// # Client-side instrumentation
//
// For outbound requests, use NewTransport to wrap a base http.RoundTripper (for
// example http.DefaultTransport) and install the returned transport on an
// http.Client.
//
// Example:
//
//	rt := telemetry.NewTransport(http.DefaultTransport)
//	c := &http.Client{Transport: rt}
//
// # Client tracing (httptrace integration)
//
// The NewClientTrace helper creates a net/http/httptrace.ClientTrace that
// integrates HTTP client request lifecycle events with OpenTelemetry.
// WithClientTrace can be used to provide a function that creates the trace from
// a request context when configuring otelhttp transport options.
//
// These helpers are thin wrappers around otelhttptrace.NewClientTrace and
// otelhttp.WithClientTrace, respectively.
//
// # Configuration and responsibilities
//
// The helper constructors accept upstream option types (otelhttp.Option and
// otelhttptrace.ClientTraceOption). Refer to the upstream documentation for the
// exact semantics and supported options for the version you vend.
//
// This package does not configure OpenTelemetry SDK state (exporters, sampling,
// resources, propagators, or global providers). Configure those elsewhere in
// your application.
package telemetry
