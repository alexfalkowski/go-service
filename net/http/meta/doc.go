// Package meta provides HTTP-specific context metadata helpers and middleware for go-service.
//
// This package serves two related purposes for HTTP request handling:
//
//   - It exposes small convenience wrappers around the generic `meta` package for exporting
//     context-scoped attributes as string maps suitable for logging and header propagation
//     (for example CamelStrings).
//
//   - It provides a small context-backed store for request-scoped HTTP objects used by go-service
//     handlers and middleware, including:
//
//   - the incoming `*http.Request`
//
//   - the active `http.ResponseWriter`
//
//   - the negotiated `encoding.Encoder` (typically selected from the request Content-Type)
//
// # Safety and expectations
//
// Request, Response, and Encoder are intentionally strict helpers: they expect the corresponding values
// to have been stored in the context via WithContent. Calling them without content metadata present will
// panic due to type assertions.
//
// These helpers are typically used in tightly controlled handler pipelines (for example those created by
// `net/http/content.NewHandler` / `NewRequestHandler`), which populate the context before invoking
// downstream logic.
//
// # Forwarded IP trust boundary
//
// HTTP server metadata extraction intentionally treats common forwarding
// headers, such as X-Forwarded-For, X-Real-IP, CF-Connecting-IP, and
// True-Client-IP, as trusted inputs and prefers them over RemoteAddr.
//
// This package does not fetch CDN provider IP ranges, maintain trusted proxy
// CIDR lists, or decide whether a request came through a trusted edge. That
// policy belongs at the infrastructure boundary: CDN, ingress, load balancer,
// firewall, service mesh, or application-specific middleware.
//
// Deployments that use the extracted IP for access logs, policy, or rate
// limiting should ensure direct origin access is blocked and the trusted edge
// strips or overwrites client-supplied forwarding headers before traffic reaches
// the service.
//
// This package also provides HTTP metadata middleware via `NewHandler` and `NewRoundTripper`.
package meta
