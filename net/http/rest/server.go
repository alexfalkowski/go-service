package rest

import (
	"fmt"
	"net/http"

	"github.com/alexfalkowski/go-service/net/http/content"
)

// Delete for rest.
func Delete[Res any](path string, handler content.Handler[Res]) {
	Route(fmt.Sprintf("%s %s", http.MethodDelete, path), handler)
}

// Get for rest.
func Get[Res any](path string, handler content.Handler[Res]) {
	Route(fmt.Sprintf("%s %s", http.MethodGet, path), handler)
}

// Post for rest.
func Post[Req any, Res any](path string, handler content.RequestHandler[Req, Res]) {
	RouteRequest(fmt.Sprintf("%s %s", http.MethodPost, path), handler)
}

// Put for rest.
func Put[Req any, Res any](path string, handler content.RequestHandler[Req, Res]) {
	RouteRequest(fmt.Sprintf("%s %s", http.MethodPut, path), handler)
}

// Patch for rest.
func Patch[Req any, Res any](path string, handler content.RequestHandler[Req, Res]) {
	RouteRequest(fmt.Sprintf("%s %s", http.MethodPatch, path), handler)
}

// RouteRequest for rest.
func RouteRequest[Req any, Res any](path string, handler content.RequestHandler[Req, Res]) {
	mux.HandleFunc(path, content.NewRequestHandler(cont, "rest", handler))
}

// Route for rest.
func Route[Res any](path string, handler content.Handler[Res]) {
	mux.HandleFunc(path, content.NewHandler(cont, "rest", handler))
}
