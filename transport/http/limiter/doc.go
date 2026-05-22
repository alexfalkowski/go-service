// Package limiter provides HTTP rate limiter middleware and wiring for go-service.
//
// This package integrates rate limiting into HTTP servers (handler middleware) and HTTP clients
// (RoundTripper middleware).
//
// Security note: treat the built-in limiter as a last-resort local protection
// layer, not as the primary production abuse-control mechanism. Prefer an
// external limiter at the edge, gateway, ingress, load balancer, or service-mesh
// boundary for consistent enforcement before requests spend application CPU.
//
// Server-side HTTP limiter middleware is installed after token verification in
// the standard transport chain. Requests with missing or invalid authorization
// are rejected before they reach this limiter.
//
// Start with `NewHandler` for server-side limiting and `NewRoundTripper` for client-side limiting.
package limiter
