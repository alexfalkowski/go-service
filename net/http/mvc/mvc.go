package mvc

import (
	"context"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/alexfalkowski/go-service/meta"
	hc "github.com/alexfalkowski/go-service/net/http/context"
	"github.com/alexfalkowski/go-service/net/http/status"
)

var mux *http.ServeMux

// Register for mvc.
func Register(mu *http.ServeMux) {
	mux = mu
}

type (
	// Result from the controller.
	Result struct {
		model any
		view  *template.Template
	}

	// Controller for mvc.
	Controller func(ctx context.Context) *Result
)

// NewResult for mvc.
func NewResult(model any, view *template.Template) *Result {
	return &Result{model: model, view: view}
}

// View from fs with path.
func View(fs fs.FS, path string) *template.Template {
	return template.Must(template.ParseFS(fs, path))
}

// Route the path with controller for mvc.
func Route(path string, controller Controller) {
	h := func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "text/html; charset=utf-8")

		ctx := req.Context()
		ctx = hc.WithRequest(ctx, req)
		ctx = hc.WithResponse(ctx, res)

		r := controller(ctx)

		if err, ok := r.model.(error); ok {
			meta.WithAttribute(ctx, "mvcModelError", meta.Error(err))
			res.WriteHeader(status.Code(err))
		}

		if err := r.view.Execute(res, r.model); err != nil {
			meta.WithAttribute(ctx, "mvcViewError", meta.Error(err))
		}
	}

	mux.HandleFunc(path, h)
}
