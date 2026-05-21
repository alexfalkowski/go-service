// Package retry provides gRPC retry interceptors and wiring for go-service.
//
// This package integrates retry behavior into gRPC client calls (for example via
// client-side interceptors) and centralizes retry-related defaults used by the
// transport stack.
//
// Backward compatibility: if no policy is passed to UnaryClientInterceptor, all unary RPCs are eligible for
// retry when they hit a retryable gRPC status. New callers that only want side-effect-safe retries should pass
// IdempotentMethods, StandardReadMethods, or another explicit policy.
//
// Start with `Config` and `UnaryClientInterceptor`.
package retry
