// Package logger provides structured logging helpers and wiring for go-service.
//
// # Overview
//
// This package builds a `log/slog` logger based on `Config` and installs it as the
// process-wide default via `slog.SetDefault`.
//
// The primary entry point is `NewLogger`, which is intended to be wired with Fx/Dig
// using `LoggerParams`.
//
// # Kinds
//
// `Config.Kind` selects the handler/exporter implementation. Supported kinds are:
//
//   - "otlp": exports logs via OpenTelemetry OTLP/HTTP and bridges `slog` to OTel
//   - "json": writes JSON logs to stdout
//   - "text": writes text logs to stdout
//   - "tint": writes colorized logs to stdout (using tint)
//
// If `Config.Kind` is unknown, `NewLogger` returns `ErrNotFound`.
//
// # Context metadata and errors
//
// `Logger.Log`/`Logger.LogAttrs` attach request/service metadata derived from the
// passed context (see `Meta`) and include a standardized `error` attribute when
// `Message.Error` is set (see `Error`).
//
// # Configuration
//
// The OTLP logger uses `Config.URL` and `Config.Headers`. Header values may be
// configured as "source strings" (for example `env:NAME`, `file:/path`, or a literal
// value) and are resolved by the telemetry header helpers used by `Config` consumers.
package logger
