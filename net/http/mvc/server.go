package mvc

import (
	"io"
	"net/http"

	"github.com/alexfalkowski/go-service/meta"
	hm "github.com/alexfalkowski/go-service/net/http/meta"
	"github.com/alexfalkowski/go-service/net/http/status"
)

// Route the path with controller for mvc.
func Route[Model any](path string, controller Controller[Model]) bool {
	if !views.IsValid() {
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

			view.Render(ctx, res, err)
		} else {
			view.Render(ctx, res, model)
		}
	}

	mux.HandleFunc(path, handler)

	return true
}

// StaticFile to be served via path.
func StaticFile(path, name string) bool {
	if !views.IsValid() {
		return false
	}

	handler := func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		f, err := views.fs.Open(name)
		if err != nil {
			meta.WithAttribute(ctx, "mvcStaticFileError", meta.Error(err))
			res.WriteHeader(status.Code(err))
		} else {
			_, err := io.Copy(res, f)
			if err != nil {
				meta.WithAttribute(ctx, "mvcStaticFileError", meta.Error(err))
				res.WriteHeader(status.Code(err))
			}
		}
	}

	mux.HandleFunc(path, handler)

	return true
}

// StaticPathValue to be served from a dedicated path value.
func StaticPathValue(path, value, prefix string) bool {
	if !views.IsValid() {
		return false
	}

	handler := func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		val := req.PathValue(value)

		f, err := views.fs.Open(prefix + "/" + val)
		if err != nil {
			meta.WithAttribute(ctx, "mvcStaticPathValueError", meta.Error(err))
			res.WriteHeader(status.Code(err))
		} else {
			_, err := io.Copy(res, f)
			if err != nil {
				meta.WithAttribute(ctx, "mvcStaticPathValueError", meta.Error(err))
				res.WriteHeader(status.Code(err))
			}
		}
	}

	mux.HandleFunc(path, handler)

	return true
}
