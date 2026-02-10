package mvc

import (
	"io"
	"path/filepath"

	"github.com/alexfalkowski/go-service/v2/mime"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/net/http/meta"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// Delete registers an HTTP DELETE route that invokes controller.
func Delete[Model any](pattern string, controller Controller[Model]) bool {
	return Route(strings.Join(strings.Space, http.MethodDelete, pattern), controller)
}

// Get registers an HTTP GET route that invokes controller.
func Get[Model any](pattern string, controller Controller[Model]) bool {
	return Route(strings.Join(strings.Space, http.MethodGet, pattern), controller)
}

// Post registers an HTTP POST route that invokes controller.
func Post[Model any](pattern string, controller Controller[Model]) bool {
	return Route(strings.Join(strings.Space, http.MethodPost, pattern), controller)
}

// Put registers an HTTP PUT route that invokes controller.
func Put[Model any](pattern string, controller Controller[Model]) bool {
	return Route(strings.Join(strings.Space, http.MethodPut, pattern), controller)
}

// Patch registers an HTTP PATCH route that invokes controller.
func Patch[Model any](pattern string, controller Controller[Model]) bool {
	return Route(strings.Join(strings.Space, http.MethodPatch, pattern), controller)
}

// Route registers a handler for pattern that invokes controller and renders the returned view.
//
// It returns false when MVC is not defined (see IsDefined).
//
// The handler sets the response Content-Type to HTML and stores the request and response writer in the
// request context (via net/http/meta) before invoking the controller.
//
// If controller returns an error, the handler writes the corresponding status code (see net/http/status.Code)
// and renders the view using the error value as the template model.
func Route[Model any](pattern string, controller Controller[Model]) bool {
	if !IsDefined() {
		return false
	}

	handler := func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set(content.TypeKey, mime.HTMLMediaType)

		ctx := req.Context()
		ctx = meta.WithRequest(ctx, req)
		ctx = meta.WithResponse(ctx, res)

		view, model, err := controller(ctx)
		if err != nil {
			meta.WithAttribute(ctx, "mvcModelError", meta.Error(err))
			res.WriteHeader(status.Code(err))

			view.Render(ctx, err)
		} else {
			view.Render(ctx, model)
		}
	}

	http.HandleFunc(mux, pattern, handler)
	return true
}

// StaticFile registers an HTTP GET route that serves the named file from the registered filesystem.
//
// It returns false when MVC is not defined (see IsDefined).
func StaticFile(pattern, name string) bool {
	if !IsDefined() {
		return false
	}

	handler := func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		if err := writeFile(name, res); err != nil {
			meta.WithAttribute(ctx, "mvcStaticFileError", meta.Error(err))
			res.WriteHeader(status.Code(err))
		}
	}

	http.HandleFunc(mux, strings.Join(strings.Space, http.MethodGet, pattern), handler)
	return true
}

// StaticPathValue registers an HTTP GET route that serves a file chosen by a path value.
//
// The file name is built as filepath.Join(prefix, req.PathValue(value)).
//
// It returns false when MVC is not defined (see IsDefined).
func StaticPathValue(pattern, value, prefix string) bool {
	if !IsDefined() {
		return false
	}

	handler := func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		name := filepath.Join(prefix, req.PathValue(value))

		if err := writeFile(name, res); err != nil {
			meta.WithAttribute(ctx, "mvcStaticPathValueError", meta.Error(err))
			res.WriteHeader(status.Code(err))
		}
	}

	http.HandleFunc(mux, strings.Join(strings.Space, http.MethodGet, pattern), handler)
	return true
}

func writeFile(name string, writer io.Writer) error {
	f, err := fileSystem.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(writer, f)
	return err
}
