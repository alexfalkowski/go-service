package rest

import (
	"net/http"

	"github.com/alexfalkowski/go-service/net/http/content"
)

var (
	mux  *http.ServeMux
	cont *content.Content
)

// Register for rest.
func Register(mu *http.ServeMux, ct *content.Content) {
	mux, cont = mu, ct
}

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
