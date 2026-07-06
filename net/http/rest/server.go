package rest

import (
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// Delete registers an HTTP DELETE handler under pattern.
//
// The effective route pattern passed to the router is a method-qualified pattern of the form:
//
//	"<METHOD> <pattern>"
//
// For example:
//
//	Delete("/health", handler) // registers "DELETE /health"
//
// This helper delegates to Route.
func Delete[Res any](pattern string, handler content.Handler[Res]) {
	Route(strings.Join(strings.Space, http.MethodDelete, pattern), handler)
}

// Get registers an HTTP GET handler under pattern.
//
// The effective route pattern passed to the router is method-qualified (see Delete for details).
// This helper delegates to Route.
func Get[Res any](pattern string, handler content.Handler[Res]) {
	Route(strings.Join(strings.Space, http.MethodGet, pattern), handler)
}

// Post registers an HTTP POST handler under pattern.
//
// The effective route pattern passed to the router is method-qualified (see Delete for details).
// This helper delegates to RouteRequest.
func Post[Req any, Res any](pattern string, handler content.RequestHandler[Req, Res]) {
	RouteRequest(strings.Join(strings.Space, http.MethodPost, pattern), handler)
}

// Put registers an HTTP PUT handler under pattern.
//
// The effective route pattern passed to the router is method-qualified (see Delete for details).
// This helper delegates to RouteRequest.
func Put[Req any, Res any](pattern string, handler content.RequestHandler[Req, Res]) {
	RouteRequest(strings.Join(strings.Space, http.MethodPut, pattern), handler)
}

// Patch registers an HTTP PATCH handler under pattern.
//
// The effective route pattern passed to the router is method-qualified (see Delete for details).
// This helper delegates to RouteRequest.
func Patch[Req any, Res any](pattern string, handler content.RequestHandler[Req, Res]) {
	RouteRequest(strings.Join(strings.Space, http.MethodPatch, pattern), handler)
}

// RouteRequest registers a handler under pattern that decodes a request and encodes a response.
//
// The handler is built using [github.com/alexfalkowski/go-service/v2/net/http/content.NewRequestHandler], which:
//   - decodes the request body using Content-Type, falling back to JSON when Content-Type is absent or unknown,
//   - decodes the request body into a newly allocated request model, and
//   - encodes the returned response model using the negotiated media type.
//
// Registration:
// The resulting handler is registered on the package-level router configured via [Register].
// [Register] must be called before RouteRequest; otherwise router/cont will be nil and this function will panic.
func RouteRequest[Req any, Res any](pattern string, handler content.RequestHandler[Req, Res]) {
	router.Handle(pattern, content.NewRequestHandler(cont, handler))
}

// Route registers a handler under pattern that encodes a response.
//
// The handler is built using [github.com/alexfalkowski/go-service/v2/net/http/content.NewHandler], which:
//   - selects an encoder based on the first Accept media type, falling back to Content-Type when
//     Accept is absent, and
//   - encodes the returned response model using the negotiated media type.
//
// Registration:
// The resulting handler is registered on the package-level router configured via Register.
// Register must be called before Route; otherwise router/cont will be nil and this function will panic.
func Route[Res any](pattern string, handler content.Handler[Res]) {
	router.Handle(pattern, content.NewHandler(cont, handler))
}
