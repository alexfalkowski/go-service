// Package rest provides REST-style HTTP handler registration and client helpers for go-service.
//
// This package is built on top of [github.com/alexfalkowski/go-service/v2/net/http/content/unary].
//
// # Server-side routing
//
// [Server] (constructed via [NewServer], typically via dependency injection) registers handlers on its
// configured router using method-qualified patterns of the form:
//
//	"<METHOD> <pattern>"
//
// For example, calling server.Get("/health", handler) registers the route pattern "GET /health".
//
// The handlers are constructed using [github.com/alexfalkowski/go-service/v2/net/http/content/unary.NewHandler]
// and [github.com/alexfalkowski/go-service/v2/net/http/content/unary.NewRequestHandler], which:
//   - decode request bodies (where applicable) from Content-Type, falling back to JSON when
//     Content-Type is absent, and rejecting it with 415 when it is unparseable, unregistered, or
//     intentionally undecodable, and
//   - encode responses using the first Accept media type, falling back to Content-Type when
//     Accept is absent.
//
// Errors are written using net/http/status helpers.
//
// # Client helpers
//
// [NewClient] wraps a shared net/http/client.Client. Configure codecs, response buffering, transport,
// timeouts, and redirect policy when constructing that client.
package rest
