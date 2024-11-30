package rest

import (
	"fmt"
	"net/http"

	"github.com/alexfalkowski/go-service/net/http/content"
)

// Delete for rest.
func Delete(path string, handler content.Handler) {
	Route(fmt.Sprintf("%s %s", http.MethodDelete, path), handler)
}

// Get for rest.
func Get(path string, handler content.Handler) {
	Route(fmt.Sprintf("%s %s", http.MethodGet, path), handler)
}

// Post for rest.
func Post(path string, handler content.Handler) {
	Route(fmt.Sprintf("%s %s", http.MethodPost, path), handler)
}

// Put for rest.
func Put(path string, handler content.Handler) {
	Route(fmt.Sprintf("%s %s", http.MethodPut, path), handler)
}

// Patch for rest.
func Patch(path string, handler content.Handler) {
	Route(fmt.Sprintf("%s %s", http.MethodPatch, path), handler)
}

// Head for rest.
func Head(path string, handler content.Handler) {
	Route(fmt.Sprintf("%s %s", http.MethodHead, path), handler)
}

// Options for rest.
func Options(path string, handler content.Handler) {
	Route(fmt.Sprintf("%s %s", http.MethodOptions, path), handler)
}

// Route for rest.
func Route(path string, handler content.Handler) {
	h := cont.NewHandler("rest", handler)

	mux.HandleFunc(path, h)
}
