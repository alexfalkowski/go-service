// Package token provides HTTP token middleware and wiring for go-service.
//
// This package integrates token-based authentication into HTTP servers (verification middleware)
// and HTTP clients (Authorization header injection).
//
// Start with `NewHandler` for server-side verification and `NewRoundTripper` for client-side injection.
package token
