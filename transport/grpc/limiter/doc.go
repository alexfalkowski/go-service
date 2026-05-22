// Package limiter provides gRPC rate limiter interceptors and wiring for go-service.
//
// This package integrates rate limiting into gRPC servers (server-side interceptors) and gRPC clients
// (client-side interceptors).
//
// Security note: treat the built-in limiter as a last-resort local protection
// layer, not as the primary production abuse-control mechanism. Prefer an
// external limiter at the edge, gateway, ingress, load balancer, or service-mesh
// boundary for consistent enforcement before requests spend application CPU.
//
// Server-side gRPC limiter interceptors are installed after metadata extraction
// and token verification in the standard transport chain. Requests with
// missing, malformed, or invalid authorization are rejected before they reach
// this limiter; that is intentional and should be handled by an external edge,
// gateway, ingress, load balancer, or service-mesh limiter when those attempts
// need quota enforcement.
//
// Start with `UnaryServerInterceptor` or `StreamServerInterceptor` for server-side limiting and
// `UnaryClientInterceptor` or `StreamClientInterceptor` for client-side limiting.
package limiter
