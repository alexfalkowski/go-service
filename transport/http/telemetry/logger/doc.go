// Package logger provides HTTP logging middleware and wiring for go-service.
//
// This package integrates request/response logging into HTTP servers (handler middleware)
// and HTTP clients (RoundTripper middleware). Ignorable paths bypass logging.
//
// Start with `NewHandler` for server-side logging and `NewRoundTripper` for client-side logging.
package logger
