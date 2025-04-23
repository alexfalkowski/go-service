package mvc

import (
	"io/fs"
	"net/http"

	"github.com/alexfalkowski/go-service/meta"
	hc "github.com/alexfalkowski/go-service/net/http/context"
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
		ctx = hc.WithRequest(ctx, req)
		ctx = hc.WithResponse(ctx, res)

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

// Static file name to be served via path.
func Static(path, name string) bool {
	if !views.IsValid() {
		return false
	}

	handler := func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		bytes, err := fs.ReadFile(views.fs, name)
		if err != nil {
			meta.WithAttribute(ctx, "mvcStaticError", meta.Error(err))
			res.WriteHeader(status.Code(err))
		} else {
			if _, err := res.Write(bytes); err != nil {
				meta.WithAttribute(ctx, "mvcStaticError", meta.Error(err))
				res.WriteHeader(status.Code(err))
			}
		}
	}

	mux.HandleFunc(path, handler)

	return true
}
