// Package retry provides gRPC retry interceptors and wiring for go-service.
//
// This package integrates retry behavior into gRPC client calls (for example via
// client-side interceptors) and centralizes retry-related defaults used by the
// transport stack.
//
// Default policy: if no policy is passed to UnaryClientInterceptor, only side-effect-safe unary RPCs are
// eligible for retry. This includes AIP-style read methods and requests carrying a request-id. In go-service,
// request-id identifies the logical request and is stable across retry attempts, so services that retry writes
// should deduplicate by request-id. Callers that need different retry eligibility can pass an explicit policy.
//
// Start with [Config] and [UnaryClientInterceptor].
package retry
