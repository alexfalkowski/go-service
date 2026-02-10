package rest

import (
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// Delete registers an HTTP DELETE handler under pattern.
func Delete[Res any](pattern string, handler content.Handler[Res]) {
	Route(strings.Join(strings.Space, http.MethodDelete, pattern), handler)
}

// Get registers an HTTP GET handler under pattern.
func Get[Res any](pattern string, handler content.Handler[Res]) {
	Route(strings.Join(strings.Space, http.MethodGet, pattern), handler)
}

// Post registers an HTTP POST handler under pattern.
func Post[Req any, Res any](pattern string, handler content.RequestHandler[Req, Res]) {
	RouteRequest(strings.Join(strings.Space, http.MethodPost, pattern), handler)
}

// Put registers an HTTP PUT handler under pattern.
func Put[Req any, Res any](pattern string, handler content.RequestHandler[Req, Res]) {
	RouteRequest(strings.Join(strings.Space, http.MethodPut, pattern), handler)
}

// Patch registers an HTTP PATCH handler under pattern.
func Patch[Req any, Res any](pattern string, handler content.RequestHandler[Req, Res]) {
	RouteRequest(strings.Join(strings.Space, http.MethodPatch, pattern), handler)
}

// RouteRequest registers a handler under pattern that decodes a request and encodes a response.
//
// The handler is built using net/http/content.NewRequestHandler and is registered on the package-level mux
// configured via Register.
func RouteRequest[Req any, Res any](pattern string, handler content.RequestHandler[Req, Res]) {
	http.HandleFunc(mux, pattern, content.NewRequestHandler(cont, handler))
}

// Route registers a handler under pattern that encodes a response.
//
// The handler is built using net/http/content.NewHandler and is registered on the package-level mux
// configured via Register.
func Route[Res any](pattern string, handler content.Handler[Res]) {
	http.HandleFunc(mux, pattern, content.NewHandler(cont, handler))
}
