// Package breaker provides gRPC circuit breaker interceptors and wiring for go-service.
//
// This package integrates circuit breaking into gRPC client calls via client-side interceptors.
// Circuit breakers are maintained per fullMethod and failures are classified by gRPC status code
// (configurable via options).
//
// Start with `UnaryClientInterceptor`.
package breaker
