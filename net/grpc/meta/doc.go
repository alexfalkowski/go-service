// Package meta provides the go-service gRPC metadata import path.
//
// The package serves two related purposes:
//
//   - it wraps common `google.golang.org/grpc/metadata` context helpers and map
//     constructors so repository code can depend on a single go-service import
//     path for gRPC metadata operations
//   - it re-exports the small subset of the root `meta` package that gRPC
//     transport code commonly needs, so callers can work with metadata maps and
//     request-scoped attributes through one package
//
// In addition, the package provides client and server interceptors that keep a
// consistent metadata contract across gRPC transports. The main keys used by
// those interceptors are:
//
//   - "user-agent"
//   - "request-id"
//   - "authorization"
//   - "geolocation"
//
// Server interceptors also emit response header metadata such as
// "service-version" and "request-id".
//
// # Request-Id semantics
//
// "request-id" identifies one logical gRPC request. Client metadata
// interceptors create it before retry interceptors run, so all retry attempts
// for the same logical request keep the same value. Server metadata
// interceptors preserve an incoming "request-id" when present, otherwise they
// generate one for the RPC before passing control to downstream handlers.
//
// Because "request-id" is stable across attempts, transports and services may
// use it as the idempotency key for retryable write operations. It is not a
// per-wire attempt id.
//
// # Forwarded IP trust boundary
//
// Server metadata extraction intentionally treats common forwarding metadata,
// such as "x-forwarded-for", "x-real-ip", "cf-connecting-ip", and
// "true-client-ip", as trusted inputs and prefers them over peer addresses.
//
// This package does not fetch CDN provider IP ranges, maintain trusted proxy
// CIDR lists, or decide whether a request came through a trusted edge. That
// policy belongs at the infrastructure boundary: CDN, ingress, load balancer,
// firewall, service mesh, or application-specific middleware.
//
// Deployments that use the extracted IP for access logs, policy, or rate
// limiting should ensure direct origin access is blocked and the trusted edge
// strips or overwrites client-supplied forwarding metadata before traffic
// reaches the service.
//
// Start with `UnaryServerInterceptor` / `StreamServerInterceptor` for
// server-side extraction and `UnaryClientInterceptor` /
// `StreamClientInterceptor` for client-side injection. Use `ExtractIncoming`
// and `ExtractOutgoing` when you need mutable copies of metadata maps.
package meta
