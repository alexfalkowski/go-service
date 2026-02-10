// Package token provides gRPC token interceptors and wiring for go-service.
//
// This package integrates token-based authentication into gRPC servers (verification interceptors)
// and gRPC clients (outgoing Authorization metadata injection).
//
// Start with `UnaryServerInterceptor` / `StreamServerInterceptor` for server-side verification and
// `UnaryClientInterceptor` / `StreamClientInterceptor` for client-side injection.
package token
