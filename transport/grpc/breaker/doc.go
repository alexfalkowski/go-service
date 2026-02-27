// Package breaker provides gRPC client-side circuit breaking for go-service.
//
// It integrates circuit breaking into gRPC clients via interceptors that guard unary RPC
// invocations with circuit breakers.
//
// # Breaker scope
//
// A separate circuit breaker is maintained per gRPC `fullMethod`, so each downstream method is
// isolated from failures in other methods.
//
// # Failure accounting
//
// The interceptor classifies whether an invocation is successful based on the gRPC status code
// of the returned error. The set of status codes treated as failures is configurable via options
// (see `WithFailureCodes` and the package defaults in `defaultOpts`).
//
// # Rejections
//
// When the breaker rejects a call (open state or too many concurrent half-open probes), the
// interceptor returns a gRPC ResourceExhausted status error.
//
// Start with `UnaryClientInterceptor`.
package breaker
