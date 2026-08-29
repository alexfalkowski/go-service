// Package logger provides HTTP logging middleware and wiring for go-service.
//
// This package integrates request/response logging into HTTP servers (handler middleware)
// and HTTP clients (RoundTripper middleware). Registered operation paths (health/metrics/etc.)
// bypass ordinary access logging on the server side but still log recovered panics.
//
// Start with [NewHandler] for server-side logging and [NewRoundTripper] for client-side logging.
package logger
