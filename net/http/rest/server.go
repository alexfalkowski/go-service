package rest

import (
	"github.com/alexfalkowski/go-service/net/http/content"
)

// Delete for rest.
func Delete(path string, handler content.Handler) {
	Route("DELETE "+path, handler)
}

// Get for rest.
func Get(path string, handler content.Handler) {
	Route("GET "+path, handler)
}

// Post for rest.
func Post(path string, handler content.Handler) {
	Route("POST "+path, handler)
}

// Put for rest.
func Put(path string, handler content.Handler) {
	Route("PUT "+path, handler)
}

// Route for rest.
func Route(path string, handler content.Handler) {
	h := cont.NewHandler("rest", handler)

	mux.HandleFunc(path, h)
}
