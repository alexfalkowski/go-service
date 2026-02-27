// Package rpc provides RPC-style HTTP handler registration and client helpers for go-service.
//
// This package is built on top of net/http/content. It relies on package-level registration (see Register)
// to supply the HTTP mux, content codec helpers, and buffer pool that are used when wiring handlers and clients.
//
// # Server-side routing
//
// Server-side helpers register POST handlers on the configured mux using method-qualified patterns of the form:
//
//	"POST <pattern>"
//
// For example, calling Route("/greet.v1.Greeter/SayHello", handler) registers the route pattern
// "POST /greet.v1.Greeter/SayHello".
//
// Handlers are constructed using net/http/content.NewRequestHandler, which:
//   - selects an encoder based on the request Content-Type,
//   - decodes the request body into a request model, and
//   - encodes the response model using the negotiated media type.
//
// Errors are written using net/http/status helpers.
//
// # Client helpers
//
// Client helpers (NewClient) build an RPC client backed by net/http/client.Client using the registered
// content and buffer pool. Redirect following is disabled by default so redirect responses are returned
// to the caller.
//
// # Registration requirement
//
// Register must be called before using any server or client helpers in this package; otherwise globals will be nil and
// handler/client construction will panic.
package rpc
