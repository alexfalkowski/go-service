// Package limiter provides gRPC rate limiter interceptors and wiring for go-service.
//
// This package integrates rate limiting into gRPC servers (server-side interceptors) and gRPC clients
// (client-side interceptors).
//
// Start with `UnaryServerInterceptor` for server-side limiting and `UnaryClientInterceptor` for client-side limiting.
package limiter
