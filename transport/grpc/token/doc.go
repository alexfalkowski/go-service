// Package token provides gRPC token interceptors and wiring for go-service.
//
// This package integrates token-based authentication into gRPC servers (verification interceptors)
// and gRPC clients (outgoing Authorization metadata injection).
//
// Shared transport access-control config constructs an injectable controller. The built-in gRPC token
// verification interceptors authenticate tokens and store the verified subject in metadata. The built-in
// access interceptors enforce configured authorization policy using that subject and the transport
// service-method.
//
// Start with [UnaryServerInterceptor] / [StreamServerInterceptor] for server-side verification and
// [UnaryAccessServerInterceptor] / [StreamAccessServerInterceptor] for server-side authorization.
// Use [UnaryClientInterceptor] / [StreamClientInterceptor] for client-side injection.
package token
