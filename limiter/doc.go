// Package limiter provides in-memory rate limiting primitives used by go-service.
//
// This package provides a configurable in-memory token-bucket style limiter that can be composed into
// transports (HTTP/gRPC) and other request-processing pipelines. Limits are applied per key derived from
// request metadata.
//
// # Keys and kinds
//
// Rate limits are keyed by a string derived from the request context. The default key derivation functions
// are provided by `NewKeyMap`, mapping a configured kind to a function that extracts a value from context
// metadata:
//
//   - "user-agent": meta.UserAgent
//   - "ip": meta.IPAddr
//   - "token": meta.Authorization
//
// The configured kind is typically provided via `Config.Kind`. If the kind is not present in the KeyMap,
// limiter construction fails with ErrMissingKey.
//
// The "ip" key uses IP metadata populated by transport metadata middleware. That
// middleware intentionally trusts common forwarding headers such as X-Forwarded-For
// and CF-Connecting-IP. Use this key only when requests first pass through trusted
// proxies that strip or overwrite client-supplied forwarding headers.
//
// # In-memory behavior and lifecycle
//
// The limiter uses an in-memory store. Limits are enforced per process and are not shared across replicas.
// This makes it suitable for single-instance deployments, development environments, or as a local
// protection mechanism, but not as a global distributed rate limit.
//
// When constructed via `NewLimiter`, the underlying store is closed on application shutdown via an Fx/Dig
// lifecycle hook.
//
// Start with `Config`, `KeyMap`/`NewKeyMap`, and `NewLimiter`.
package limiter
