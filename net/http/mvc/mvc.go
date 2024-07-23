package mvc

import (
	"context"
	"html/template"
	"net/http"

	"github.com/alexfalkowski/go-service/meta"
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
	Controller func(ctx context.Context, req *http.Request, res http.ResponseWriter) *Result
)

// NewResult for mvc.
func NewResult(model any, view *template.Template) *Result {
	return &Result{model: model, view: view}
}

// Route the path with controller for mvc.
func Route(path string, controller Controller) {
	h := func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "text/html; charset=utf-8")

		ctx := req.Context()
		r := controller(ctx, req, res)

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
