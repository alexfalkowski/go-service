// Package rest provides REST-style HTTP handler registration and client helpers for go-service.
//
// This package is built on top of net/http/content. It relies on package-level registration (see Register)
// to supply the HTTP mux, content codec helpers, and buffer pool that are used when wiring handlers and clients.
//
// # Server-side routing
//
// Server-side helpers (Get/Post/etc.) register handlers on the configured mux using method-qualified
// patterns of the form:
//
//	"<METHOD> <pattern>"
//
// For example, calling Get("/health", handler) registers the route pattern "GET /health".
//
// The handlers are constructed using net/http/content helpers (NewHandler/NewRequestHandler), which:
//   - select an encoder based on the request Content-Type,
//   - decode request bodies (where applicable), and
//   - encode responses using the negotiated media type.
//
// Errors are written using net/http/status helpers.
//
// # Client helpers
//
// Client helpers (NewClient) build a net/http/client.Client using the registered content and buffer pool.
// Redirect following is disabled by default so redirect responses are returned to the caller.
//
// # Registration requirement
//
// Register must be called before using any server or client helpers in this package; otherwise globals will be nil and
// handler/client construction will panic.
package rest
