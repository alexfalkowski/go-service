// Package limiter provides HTTP rate limiter middleware and wiring for go-service.
//
// This package integrates rate limiting into HTTP servers (handler middleware) and HTTP clients
// (RoundTripper middleware).
//
// Start with `NewHandler` for server-side limiting and `NewRoundTripper` for client-side limiting.
package limiter
