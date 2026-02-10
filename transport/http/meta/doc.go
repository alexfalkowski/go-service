// Package meta provides HTTP metadata middleware and wiring for go-service.
//
// This package extracts request metadata into the request context on the server side and injects
// outbound request metadata on the client side.
//
// Start with `NewHandler` for server-side extraction and `NewRoundTripper` for client-side injection.
package meta
