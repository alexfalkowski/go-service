// Package limiter provides gRPC rate limiter interceptors and wiring for go-service.
//
// This package integrates rate limiting into gRPC servers (server-side interceptors) and gRPC clients
// (client-side interceptors).
//
// Start with `UnaryServerInterceptor` or `StreamServerInterceptor` for server-side limiting and
// `UnaryClientInterceptor` or `StreamClientInterceptor` for client-side limiting.
package limiter
