package rest

import (
	"fmt"
	"net/http"

	"github.com/alexfalkowski/go-service/net/http/content"
)

// Delete for rest.
func Delete[Req any, Res any](path string, handler content.Handler[Req, Res]) {
	RouteQuery(fmt.Sprintf("%s %s", http.MethodDelete, path), handler)
}

// Get for rest.
func Get[Req any, Res any](path string, handler content.Handler[Req, Res]) {
	RouteQuery(fmt.Sprintf("%s %s", http.MethodGet, path), handler)
}

// Post for rest.
func Post[Req any, Res any](path string, handler content.Handler[Req, Res]) {
	RouteBody(fmt.Sprintf("%s %s", http.MethodPost, path), handler)
}

// Put for rest.
func Put[Req any, Res any](path string, handler content.Handler[Req, Res]) {
	RouteBody(fmt.Sprintf("%s %s", http.MethodPut, path), handler)
}

// Patch for rest.
func Patch[Req any, Res any](path string, handler content.Handler[Req, Res]) {
	RouteBody(fmt.Sprintf("%s %s", http.MethodPatch, path), handler)
}

// RouteBody for rest.
func RouteBody[Req any, Res any](path string, handler content.Handler[Req, Res]) {
	mux.HandleFunc(path, content.NewBodyHandler(cont, "rest", handler))
}

// RouteQuery for rest.
func RouteQuery[Req any, Res any](path string, handler content.Handler[Req, Res]) {
	mux.HandleFunc(path, content.NewQueryHandler(cont, "rest", handler))
}
