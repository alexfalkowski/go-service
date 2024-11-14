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
	var t *template.Template

	if params.FS == nil || params.Patterns == nil {
		t = template.New("")
	} else {
		t = template.Must(template.New("").Funcs(sprigin.FuncMap()).ParseFS(params.FS, params.Patterns...))
	}

	return &Views{template: t, fs: params.FS}
}

// NewRouter for mvc.
func NewRouter(mux *http.ServeMux, views *Views) *Router {
	return &Router{mux: mux, views: views}
}

// Route the path with controller for mvc.
func (r *Router) Route(path string, controller Controller) {
	h := func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "text/html; charset=utf-8")

		ctx := req.Context()
		ctx = hc.WithRequest(ctx, req)
		ctx = hc.WithResponse(ctx, res)

		v, m := controller(ctx)

		if err, ok := m.(error); ok {
			meta.WithAttribute(ctx, "mvcModelError", meta.Error(err))
			res.WriteHeader(status.Code(err))
		}

		if err := r.views.template.ExecuteTemplate(res, string(v), m); err != nil {
			meta.WithAttribute(ctx, "mvcViewError", meta.Error(err))
		}
	}

	r.mux.HandleFunc(path, h)
}

// Static file name to be served via path.
func (r *Router) Static(path, name string) {
	fs := r.views.fs
	if fs == nil {
		return
	}

	h := func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		b, err := fs.ReadFile(name)
		if err != nil {
			meta.WithAttribute(ctx, "mvcStaticError", meta.Error(err))
			res.WriteHeader(status.Code(err))

			return
		}

		if _, err := res.Write(b); err != nil {
			meta.WithAttribute(ctx, "mvcStaticError", meta.Error(err))

			return
		}
	}

	r.mux.HandleFunc(path, h)
}
