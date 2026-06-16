// Package token provides gRPC token interceptors and wiring for go-service.
//
// This package integrates token-based authentication into gRPC servers (verification interceptors)
// and gRPC clients (outgoing Authorization metadata injection).
//
// Access-control config constructs an injectable [AccessController]. The built-in gRPC token
// interceptors authenticate tokens and store the verified subject in metadata; they do not enforce
// authorization policy automatically. Services that need authorization must call [AccessController.HasAccess]
// from handlers or install their own interceptor using the verified user id.
//
// Start with [UnaryServerInterceptor] / [StreamServerInterceptor] for server-side verification and
// [UnaryClientInterceptor] / [StreamClientInterceptor] for client-side injection.
package token
