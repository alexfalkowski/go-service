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

// Delete method for mvc.
func Delete[Model any](pattern string, controller Controller[Model]) bool {
	return Route(strings.Join(strings.Space, http.MethodDelete, pattern), controller)
}

// Get method for mvc.
func Get[Model any](pattern string, controller Controller[Model]) bool {
	return Route(strings.Join(strings.Space, http.MethodGet, pattern), controller)
}

// Post method for mvc.
func Post[Model any](pattern string, controller Controller[Model]) bool {
	return Route(strings.Join(strings.Space, http.MethodPost, pattern), controller)
}

// Put method for mvc.
func Put[Model any](pattern string, controller Controller[Model]) bool {
	return Route(strings.Join(strings.Space, http.MethodPut, pattern), controller)
}

// Patch method for mvc.
func Patch[Model any](pattern string, controller Controller[Model]) bool {
	return Route(strings.Join(strings.Space, http.MethodPatch, pattern), controller)
}

// Route the path with controller for mvc.
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

	mux.HandleFunc(pattern, handler)
	return true
}

// StaticFile to be served via path.
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

	mux.HandleFunc(strings.Join(strings.Space, http.MethodGet, pattern), handler)
	return true
}

// StaticPathValue to be served from a dedicated path value.
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

	mux.HandleFunc(strings.Join(strings.Space, http.MethodGet, pattern), handler)
	return true
}

func writeFile(name string, writer io.Writer) error {
	f, err := fileSystem.Open(name)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, f)
	return err
}
