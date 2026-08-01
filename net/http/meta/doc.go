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
// # Safety and expectations
//
// Request and Response return nil when request-response metadata is absent. Handler pipelines that need
// those values should install them with WithRequestResponse before invoking downstream logic.
//
// These helpers are typically used in tightly controlled handler pipelines (for example those created by
// [github.com/alexfalkowski/go-service/v2/net/http/content.NewHandler] /
// [github.com/alexfalkowski/go-service/v2/net/http/content.NewRequestHandler]), which populate the context before invoking
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
// # Request-Id semantics
//
// Request-Id identifies one logical HTTP request. Client metadata middleware
// creates it before retry middleware runs, so all retry attempts for the same
// logical request keep the same value. Server metadata middleware preserves an
// incoming Request-Id when present, otherwise it generates one for the request
// before passing control to downstream handlers.
//
// Because Request-Id is stable across attempts, transports and services may use
// it as the idempotency key for retryable write operations. It is not a per-wire
// attempt id.
//
// This package also provides HTTP metadata middleware via [NewHandler] and [NewRoundTripper].
package meta
