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
//   - "user-id": [meta.UserID]
//   - "transport-service-method": [meta.TransportServiceMethod]
//   - "service-method": [meta.ServiceMethod]
//   - "ip": [meta.IPAddr]
//   - "user-agent": [meta.UserAgent]
//
// The configured kind is typically provided via [Config.Kind]. If the kind is not present in the KeyMap,
// limiter construction fails with ErrMissingKey.
//
// The "user-id" key uses the verified principal stored in metadata. For
// JWT/PASETO tokens this is the subject claim; for SSH tokens this is the
// verified key name returned by the token verifier. Prefer it for per-principal
// quotas when authenticated identity is available.
//
// The "transport-service-method" key prefixes that service-method value with
// the transport name, such as "http:GET /users/{id}" or
// "grpc:/users.v1.Users/Get". Use it when HTTP and gRPC operations should have
// independent limiter buckets even if their service-method values overlap.
//
// The "service-method" key uses transport metadata populated by the HTTP and
// gRPC metadata layers. HTTP stores the matched request pattern when available
// and otherwise the URL path; gRPC stores the full method name. Prefer
// "transport-service-method" unless cross-transport operations intentionally
// share quota.
//
// The "ip" key uses IP metadata populated by transport metadata middleware. That
// middleware intentionally trusts common forwarding headers such as
// X-Forwarded-For and CF-Connecting-IP. The limiter does not validate CDN or
// proxy source ranges itself; it consumes the metadata provided by the transport
// stack. Use it only when the service is reachable through trusted edge
// infrastructure that strips or overwrites client-supplied forwarding headers.
//
// The "user-agent" key uses caller-supplied User-Agent metadata. Treat it as a
// coarse development or service-to-service convenience key, not as a public-edge
// abuse-control identity.
//
// # In-memory behavior and lifecycle
//
// The limiter uses an in-memory store. Limits are enforced per process and are not shared across replicas.
// This makes it suitable for single-instance deployments, development environments, or as a last-resort
// local safeguard, but not as a global distributed rate limit.
//
// The limiter caps the number of caller-derived keys that receive independent in-memory buckets. Once
// [Config.MaxKeys] is reached, additional distinct keys share one overflow bucket so high-cardinality
// key floods cannot grow store memory without bound.
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
