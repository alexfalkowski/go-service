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
// Server-side HTTP limiter middleware is installed after metadata extraction
// and token verification in the standard transport chain. Requests with
// missing, malformed, or invalid authorization are rejected before they reach
// this limiter; that is intentional and should be handled by an external edge,
// gateway, ingress, load balancer, or service-mesh limiter when those attempts
// need quota enforcement.
//
// Start with [NewHandler] for server-side limiting and [NewRoundTripper] for client-side limiting.
package limiter
