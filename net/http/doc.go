// Package http provides small HTTP wrappers and helpers around the standard library net/http package.
//
// This package primarily re-exports common net/http types and constants behind go-service aliases and
// provides a few convenience helpers used by transport wiring, such as:
//
//   - NewClient, which wraps a RoundTripper with OpenTelemetry instrumentation and applies a client timeout,
//   - NewServer, which builds an http.Server using configured timeouts and protocol settings,
//   - Handle/HandleFunc, which register handlers wrapped with OpenTelemetry instrumentation,
//   - Pattern and ParseServiceMethod, which help standardize route naming for telemetry.
//
// Start with `NewClient` and `NewServer`.
package http
