// Package retry provides gRPC retry interceptors and wiring for go-service.
//
// This package integrates retry behavior into gRPC client calls (for example via
// client-side interceptors) and centralizes retry-related defaults used by the
// transport stack.
//
// Start with `Config` and `UnaryClientInterceptor`.
package retry
