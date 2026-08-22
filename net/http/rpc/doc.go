// Package rpc provides RPC-style HTTP handler registration and client helpers for go-service.
//
// This package is built on top of [github.com/alexfalkowski/go-service/v2/net/http/content/unary].
//
// # Server-side routing
//
// [Server] (constructed via [NewServer], typically via dependency injection) registers POST handlers on
// its configured router using method-qualified patterns of the form:
//
//	"POST <pattern>"
//
// For example, calling server.Route("/greet.v1.Greeter/SayHello", handler) registers the route pattern
// "POST /greet.v1.Greeter/SayHello".
//
// Handlers are constructed using [github.com/alexfalkowski/go-service/v2/net/http/content/unary.NewRequestHandler], which:
//   - decodes the request body into a request model from Content-Type, falling back to JSON when
//     Content-Type is absent, and rejecting it with 415 when it is unparseable, unregistered, or
//     intentionally undecodable, and
//   - encodes the response model using the first Accept media type, falling back to Content-Type when
//     Accept is absent.
//
// Errors are written using net/http/status helpers.
//
// # Client helpers
//
// [NewClient] wraps a shared net/http/client.Client. Configure codecs, response buffering, transport,
// timeouts, and redirect policy when constructing that client.
package rpc
