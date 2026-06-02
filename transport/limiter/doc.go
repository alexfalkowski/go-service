// Package limiter provides in-memory rate limiting primitives used by go-service.
//
// This package provides a configurable in-memory token-bucket style limiter that can be composed into
// transports (HTTP/gRPC) and other request-processing pipelines. Limits are applied per key derived from
// request metadata.
//
// # Keys and kinds
//
// Rate limits are keyed by a string derived from the request context. The default key derivation functions
// are provided by [NewKeyMap], mapping a configured kind to a function that extracts a value from context
// metadata:
//
//   - "user-agent": [meta.UserAgent]
//   - "ip": [meta.IPAddr]
//   - "user-id": [meta.UserID]
//   - "service-method": [meta.ServiceMethod]
//
// The configured kind is typically provided via [Config.Kind]. If the kind is not present in the KeyMap,
// limiter construction fails with ErrMissingKey.
//
// The "ip" key uses IP metadata populated by transport metadata middleware. That
// middleware intentionally trusts common forwarding headers such as
// X-Forwarded-For and CF-Connecting-IP. The limiter does not validate CDN or
// proxy source ranges itself; it consumes the metadata provided by the transport
// stack.
//
// Use the "ip" key when the service is only reachable through trusted edge
// infrastructure that strips or overwrites client-supplied forwarding headers.
// If direct origin access is possible, prefer a verified identity key such as
// "user-id" or add application-specific trusted proxy validation before relying
// on IP-based limits.
//
// The "user-id" key uses the verified principal stored in metadata. For
// JWT/PASETO tokens this is the subject claim; for SSH tokens this is the
// verified key name returned by the token verifier.
//
// The "service-method" key uses transport metadata populated by the HTTP and
// gRPC metadata layers. HTTP stores the matched request pattern when available
// and otherwise the URL path; gRPC stores the full method name.
//
// # In-memory behavior and lifecycle
//
// The limiter uses an in-memory store. Limits are enforced per process and are not shared across replicas.
// This makes it suitable for single-instance deployments, development environments, or as a last-resort
// local safeguard, but not as a global distributed rate limit.
//
// For primary production abuse protection, prefer an external rate limiter at
// the edge, gateway, ingress, load balancer, or service-mesh boundary. Those
// layers can enforce limits consistently across replicas and before requests
// spend application CPU.
//
// When constructed via [NewLimiter], the underlying store is closed on application shutdown via an Fx/Dig
// lifecycle hook.
//
// Start with [Config], [KeyMap]/[NewKeyMap], and [NewLimiter].
package limiter
