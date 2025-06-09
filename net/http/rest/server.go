package rest

import (
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// Delete for rest.
func Delete[Res any](pattern string, handler content.Handler[Res]) {
	Route(strings.Join(" ", http.MethodDelete, pattern), handler)
}

// Get for rest.
func Get[Res any](pattern string, handler content.Handler[Res]) {
	Route(strings.Join(" ", http.MethodGet, pattern), handler)
}

// Post for rest.
func Post[Req any, Res any](pattern string, handler content.RequestHandler[Req, Res]) {
	RouteRequest(strings.Join(" ", http.MethodPost, pattern), handler)
}

// Put for rest.
func Put[Req any, Res any](pattern string, handler content.RequestHandler[Req, Res]) {
	RouteRequest(strings.Join(" ", http.MethodPut, pattern), handler)
}

// Patch for rest.
func Patch[Req any, Res any](pattern string, handler content.RequestHandler[Req, Res]) {
	RouteRequest(strings.Join(" ", http.MethodPatch, pattern), handler)
}

// RouteRequest for rest.
func RouteRequest[Req any, Res any](pattern string, handler content.RequestHandler[Req, Res]) {
	mux.HandleFunc(pattern, content.NewRequestHandler(cont, handler))
}

// Route for rest.
func Route[Res any](pattern string, handler content.Handler[Res]) {
	mux.HandleFunc(pattern, content.NewHandler(cont, handler))
}
