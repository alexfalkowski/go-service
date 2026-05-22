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
// Server-side gRPC limiter interceptors are installed after token verification
// in the standard transport chain. Requests with missing or invalid
// authorization are rejected before they reach this limiter.
//
// Start with `UnaryServerInterceptor` or `StreamServerInterceptor` for server-side limiting and
// `UnaryClientInterceptor` or `StreamClientInterceptor` for client-side limiting.
package limiter
