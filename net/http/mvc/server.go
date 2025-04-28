package mvc

import (
	"io"
	"net/http"

	"github.com/alexfalkowski/go-service/meta"
	hm "github.com/alexfalkowski/go-service/net/http/meta"
	"github.com/alexfalkowski/go-service/net/http/status"
)

// Delete method for mvc.
func Delete[Model any](path string, controller Controller[Model]) bool {
	return Route(http.MethodDelete+" "+path, controller)
}

// Get method for mvc.
func Get[Model any](path string, controller Controller[Model]) bool {
	return Route(http.MethodGet+" "+path, controller)
}

// Post method for mvc.
func Post[Model any](path string, controller Controller[Model]) bool {
	return Route(http.MethodPost+" "+path, controller)
}

// Put method for mvc.
func Put[Model any](path string, controller Controller[Model]) bool {
	return Route(http.MethodPut+" "+path, controller)
}

// Patch method for mvc.
func Patch[Model any](path string, controller Controller[Model]) bool {
	return Route(http.MethodPatch+" "+path, controller)
}

// Route the path with controller for mvc.
func Route[Model any](path string, controller Controller[Model]) bool {
	if !IsDefined() {
		return false
	}

	handler := func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "text/html; charset=utf-8")

		ctx := req.Context()
		ctx = hm.WithRequest(ctx, req)
		ctx = hm.WithResponse(ctx, res)

		view, model, err := controller(ctx)
		if err != nil {
			meta.WithAttribute(ctx, "mvcModelError", meta.Error(err))
			res.WriteHeader(status.Code(err))

			view.Render(ctx, err)
		} else {
			view.Render(ctx, model)
		}
	}

	mux.HandleFunc(path, handler)

	return true
}

// StaticFile to be served via path.
func StaticFile(path, name string) bool {
	if !IsDefined() {
		return false
	}

	path = http.MethodGet + " " + path
	handler := func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		if err := writeFile(name, res); err != nil {
			meta.WithAttribute(ctx, "mvcStaticFileError", meta.Error(err))
			res.WriteHeader(status.Code(err))
		}
	}

	mux.HandleFunc(path, handler)

	return true
}

// StaticPathValue to be served from a dedicated path value.
func StaticPathValue(path, value, prefix string) bool {
	if !IsDefined() {
		return false
	}

	path = http.MethodGet + " " + path
	handler := func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		name := prefix + "/" + req.PathValue(value)

		if err := writeFile(name, res); err != nil {
			meta.WithAttribute(ctx, "mvcStaticPathValueError", meta.Error(err))
			res.WriteHeader(status.Code(err))
		}
	}

	mux.HandleFunc(path, handler)

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
