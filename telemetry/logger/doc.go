// Package logger provides structured logging helpers and wiring for go-service.
//
// This package constructs a [log/slog] logger according to [Config], installs it as the
// process-wide default (via [slog.SetDefault]), and provides a thin wrapper (*[Logger])
// that standardizes how context metadata and errors are attached to log records.
//
// # Overview
//
// The primary entry point is [NewLogger], which is intended to be wired with Fx/Dig
// using [LoggerParams]. When enabled, NewLogger:
//
//   - builds the configured slog logger,
//   - installs it as the default logger for the process, and
//   - returns a *[Logger] wrapper that provides convenience methods.
//
// When logging is disabled (Config is nil), NewLogger returns (nil, nil).
//
// # Kinds
//
// [Config.Kind] selects the handler/exporter implementation. This package typically
// supports the following kinds:
//
//   - "otlp": exports logs via OpenTelemetry OTLP and bridges slog to OTel
//   - "json": writes JSON logs to stdout
//   - "text": writes text logs to stdout
//   - "tint": writes colorized logs to stdout (using tint)
//
// If [Config.Kind] is unknown, NewLogger returns [ErrNotFound].
//
// # Context metadata and errors
//
// [Logger.Log] / [Logger.LogAttrs] are opinionated helpers that keep log records consistent.
// They:
//
//   - append attributes derived from the provided context (via [Meta]), and
//   - append a standardized "error" attribute when [Message.Error] is non-nil (via [Error]).
//
// This ensures common request/service metadata and errors show up consistently across
// handlers/exporters.
//
// Metadata values are capped at 1024 bytes before being attached to the log
// record. Truncation preserves UTF-8 validity so multi-byte characters are not
// split.
//
// # Exporter headers and secret resolution
//
// Some kinds (notably "otlp") support outbound request headers for authentication or
// routing. These are configured in [Config.Headers].
//
// Header values may be configured using go-service "source strings" (for example
// "env:NAME", "file:/path", or a literal value). Those values are resolved by
// [github.com/alexfalkowski/go-service/v2/telemetry/header.Map.Secrets] or [github.com/alexfalkowski/go-service/v2/telemetry/header.Map.MustSecrets] by the consumer
// that projects configuration before constructing exporters.
//
// When [Config.Headers] is non-empty, non-loopback "http://" OTLP endpoints are
// rejected to avoid sending credential-bearing headers over cleartext transport.
// Use "https://" for external collectors. Local development collectors on
// "localhost" or loopback IP addresses may use "http://".
// Header-bearing remote OTLP/gRPC endpoints are rejected until TLS configuration
// is supported for OTLP/gRPC exporters.
//
// # Notes
//
// This package is primarily wiring and thin adaptation. For exact exporter semantics
// and supported configuration options, consult the upstream OpenTelemetry components
// used by the selected kind.
package logger
