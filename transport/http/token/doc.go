// Package token provides HTTP token middleware and wiring for go-service.
//
// This package integrates token-based authentication into HTTP servers (verification middleware),
// token-based authorization into HTTP servers (access-control middleware), and HTTP clients
// (Authorization header injection).
//
// Start with [NewHandler] for server-side verification, [NewAccessHandler] for server-side access control,
// and [NewRoundTripper] for client-side injection.
package token
