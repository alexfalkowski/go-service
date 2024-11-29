package rest

import (
	"fmt"

	"github.com/alexfalkowski/go-service/net/http/content"
	"github.com/go-resty/resty/v2"
)

// Delete for rest.
func Delete(path string, handler content.Handler) {
	Route(fmt.Sprintf("%s %s", resty.MethodDelete, path), handler)
}

// Get for rest.
func Get(path string, handler content.Handler) {
	Route(fmt.Sprintf("%s %s", resty.MethodGet, path), handler)
}

// Post for rest.
func Post(path string, handler content.Handler) {
	Route(fmt.Sprintf("%s %s", resty.MethodPost, path), handler)
}

// Put for rest.
func Put(path string, handler content.Handler) {
	Route(fmt.Sprintf("%s %s", resty.MethodPut, path), handler)
}

// Patch for rest.
func Patch(path string, handler content.Handler) {
	Route(fmt.Sprintf("%s %s", resty.MethodPatch, path), handler)
}

// Head for rest.
func Head(path string, handler content.Handler) {
	Route(fmt.Sprintf("%s %s", resty.MethodHead, path), handler)
}

// Options for rest.
func Options(path string, handler content.Handler) {
	Route(fmt.Sprintf("%s %s", resty.MethodOptions, path), handler)
}

// Route for rest.
func Route(path string, handler content.Handler) {
	h := cont.NewHandler("rest", handler)

	mux.HandleFunc(path, h)
}
