package mvc

import (
	"context"
	"embed"
	"html/template"
	"net/http"

	"github.com/alexfalkowski/go-service/meta"
	hc "github.com/alexfalkowski/go-service/net/http/context"
	"github.com/alexfalkowski/go-service/net/http/status"
	"github.com/go-sprout/sprout/sprigin"
	"go.uber.org/fx"
)

type (
	// ViewsParams for mvc.
	ViewsParams struct {
		fx.In

		FS       *embed.FS `optional:"true"`
		Patterns Patterns  `optional:"true"`
	}

	// View for mvc.
	Views struct {
		template *template.Template
		fs       *embed.FS
	}

	// View to render.
	View string

	// Model for mvc.
	Model any

	// Patterns to render views.
	Patterns []string

	// Router for mvc.
	Router struct {
		mux   *http.ServeMux
		views *Views
	}

	// Controller for mvc.
	Controller func(ctx context.Context) (View, Model)
)

// NewView from fs with patterns.
func NewViews(params ViewsParams) *Views {
	var tpl *template.Template

	if params.FS == nil || params.Patterns == nil {
		tpl = template.New("")
	} else {
		tpl = template.Must(template.New("").Funcs(sprigin.FuncMap()).ParseFS(params.FS, params.Patterns...))
	}

	return &Views{template: tpl, fs: params.FS}
}

// NewRouter for mvc.
func NewRouter(mux *http.ServeMux, views *Views) *Router {
	return &Router{mux: mux, views: views}
}

// Route the path with controller for mvc.
func (r *Router) Route(path string, controller Controller) {
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
		}
	}

	r.mux.HandleFunc(path, handler)
}

// Static file name to be served via path.
func (r *Router) Static(path, name string) {
	fs := r.views.fs
	if fs == nil {
		return
	}

	handler := func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		bytes, err := fs.ReadFile(name)
		if err != nil {
			meta.WithAttribute(ctx, "mvcStaticError", meta.Error(err))
			res.WriteHeader(status.Code(err))

			return
		}

		if _, err := res.Write(bytes); err != nil {
			meta.WithAttribute(ctx, "mvcStaticError", meta.Error(err))

			return
		}
	}

	r.mux.HandleFunc(path, handler)
}
