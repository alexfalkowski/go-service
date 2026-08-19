// Package http provides small HTTP wrappers and helpers around the standard library net/http package.
//
// This package primarily re-exports common net/http types and constants behind go-service aliases and
// provides a few convenience helpers used by transport wiring, such as:
//
//   - NewClient, which applies a client timeout and wraps a RoundTripper with OpenTelemetry instrumentation
//     when tracing or metrics are enabled,
//   - NewServer, which builds an [http.Server] using configured timeouts and protocol settings,
//   - NewTelemetryHandler, which wraps a handler with OpenTelemetry server instrumentation when tracing or
//     metrics are enabled,
//   - Router, which registers handlers on a mux together with route policy used by transport middleware,
//   - Pattern and ParseServiceMethod, which help standardize route naming for telemetry.
//   - AcceptItems, FirstAcceptItem, IsAcceptZeroQuality, and IsAcceptWildcard, which support Accept header handling.
//
// Server construction reads timeout keys from options.Map (`read_timeout`, `write_timeout`,
// `idle_timeout`, `read_header_timeout`), each defaulting independently to 30 seconds, and also
// supports `max_header_bytes` as an SI size string. HTTP/2 tuning is opt-in through
// `http2_max_concurrent_streams`, `http2_max_receive_buffer_per_connection`, and
// `http2_max_receive_buffer_per_stream`; unset options preserve the Go HTTP/2 defaults.
//
// Start with [NewClient] and [NewServer].
package http
