// Package breaker provides HTTP client-side circuit breaking for go-service.
//
// It integrates circuit breaking into HTTP clients by wrapping an `http.RoundTripper`.
//
// # Breaker scope
//
// Breakers are maintained per request key (method + host), so each upstream is isolated.
// Breaker keys are retained for the lifetime of the RoundTripper. This is intended for service-to-service
// clients with a small, bounded set of configured upstream hosts. Avoid enabling this wrapper for clients that
// call arbitrary user-supplied or otherwise high-cardinality hosts unless the caller bounds that host set.
//
// # Failure accounting vs caller behavior
//
// The wrapped `RoundTripper` classifies outcomes for breaker accounting:
//   - Transport errors (a non-nil error returned by the underlying `RoundTripper`) are counted as failures.
//   - HTTP responses with status codes classified as failures (see `WithFailureStatusFunc` / `WithFailureStatuses`)
//     are also counted as failures.
//
// When a response status code is treated as a failure for breaker accounting, the wrapper still returns the
// original `*http.Response` to the caller with a nil error. This allows application logic to continue to
// handle HTTP responses normally, while the circuit breaker learns about upstream health.
//
// Start with `NewRoundTripper`.
package breaker
