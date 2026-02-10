// Package limiter provides rate limiting primitives used by go-service.
//
// This package provides a configurable in-memory limiter implementation that can be composed into transports.
// Limits are applied per key derived from request metadata (for example user-agent, ip, or token).
//
// Start with `Config`, `NewKeyMap`, and `NewLimiter`.
package limiter
