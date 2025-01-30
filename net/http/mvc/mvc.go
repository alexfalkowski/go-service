package mvc

import (
	"context"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/alexfalkowski/go-service/meta"
	nh "github.com/alexfalkowski/go-service/net/http"
	hc "github.com/alexfalkowski/go-service/net/http/context"
	"github.com/alexfalkowski/go-service/net/http/status"
	"github.com/go-sprout/sprout/sprigin"
	"go.uber.org/fx"
)

type (
	// ViewsParams for mvc.
	ViewsParams struct {
		fx.In

		FS       fs.FS    `optional:"true"`
		Patterns Patterns `optional:"true"`
	}

	// Patterns to render views.
	Patterns []string
)

// IsValid verifies the params are present.
func (p ViewsParams) IsValid() bool {
	return p.FS != nil && len(p.Patterns) != 0
}

// NewView from fs with patterns.
func NewViews(params ViewsParams) *Views {
	var tpl *template.Template

	if params.IsValid() {
		tpl = template.Must(template.New("").Funcs(sprigin.FuncMap()).ParseFS(params.FS, params.Patterns...))
	}

	return &Views{template: tpl, fs: params.FS}
}

// View for mvc.
type Views struct {
	template *template.Template
	fs       fs.FS
}

// IsValid verifies that ut has an fs and template.
func (v *Views) IsValid() bool {
	return v.template != nil && v.fs != nil
}

// NewRouter for mvc.
func NewRouter(mux *http.ServeMux, views *Views) *Router {
	return &Router{mux: mux, views: views}
}

type (
	// Router for mvc.
	Router struct {
		mux   *http.ServeMux
		views *Views
	}

	// View to render.
	View string

	// Model for mvc.
	Model any

	// Controller for mvc.
	Controller func(ctx context.Context) (View, Model)
)

// Route the path with controller for mvc.
func (r *Router) Route(path string, controller Controller) bool {
	if !r.views.IsValid() {
		return false
	}

	handler := func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "text/html; charset=utf-8")

		ctx := req.Context()
		ctx = hc.WithRequest(ctx, req)
		ctx = hc.WithResponse(ctx, res)

		view, model := controller(ctx)

		if err, ok := model.(error); ok {
			meta.WithAttribute(ctx, "mvcModelError", meta.Error(err))
			res.WriteHeader(status.Code(err))
		}

		if err := r.views.template.ExecuteTemplate(res, string(view), model); err != nil {
			meta.WithAttribute(ctx, "mvcViewError", meta.Error(err))
			res.WriteHeader(status.Code(err))
		}
	}

	r.mux.HandleFunc(path, handler)

	return true
}

// Static file name to be served via path.
func (r *Router) Static(path, name string) bool {
	if !r.views.IsValid() {
		return false
	}

	handler := func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		bytes, err := fs.ReadFile(r.views.fs, name)
		if err != nil {
			meta.WithAttribute(ctx, "mvcStaticError", meta.Error(err))
			res.WriteHeader(status.Code(err))
		}

		nh.WriteResponse(ctx, res, bytes)
	}

	r.mux.HandleFunc(path, handler)

	return true
}
