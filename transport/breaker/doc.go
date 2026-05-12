// Package breaker provides circuit breaker helpers and defaults used by go-service.
//
// This package wraps and re-exports types from github.com/sony/gobreaker under a stable module
// path, and centralizes shared defaults (see DefaultSettings) used by transport integrations.
//
// Typical usage is to start from DefaultSettings, then customize:
//   - Settings.Name to a stable per-upstream/per-method key, and
//   - Settings.IsSuccessful / Settings.ReadyToTrip to control what counts as a failure and when to trip.
//
// go-service transport packages (for example the HTTP RoundTripper wrapper and the gRPC unary
// client interceptor) compose these settings to build per-destination circuit breakers.
package breaker
