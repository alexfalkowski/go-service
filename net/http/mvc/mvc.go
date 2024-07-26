package mvc

import (
	"context"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/alexfalkowski/go-service/meta"
	hc "github.com/alexfalkowski/go-service/net/http/context"
	"github.com/alexfalkowski/go-service/net/http/status"
	"github.com/go-sprout/sprout"
)

var mux *http.ServeMux

// Register for mvc.
func Register(mu *http.ServeMux) {
	mux = mu
}

type (
	// View for mvc.
	View struct {
		template *template.Template
		name     string
	}

	// Model for mvc.
	Model any

	// Controller for mvc.
	Controller func(ctx context.Context) (*View, Model)
)

// View from fs with path.
func NewView(fs fs.FS, path string) *View {
	d, f := filepath.Split(path)
	t := template.Must(template.New(d).Funcs(sprout.FuncMap(sprout.WithLogger(nil))).ParseFS(fs, path))

	return &View{name: f, template: t}
}

// Route the path with controller for mvc.
func Route(path string, controller Controller) {
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

		if err := v.template.ExecuteTemplate(res, v.name, m); err != nil {
			meta.WithAttribute(ctx, "mvcViewError", meta.Error(err))
		}
	}

	mux.HandleFunc(path, h)
}
