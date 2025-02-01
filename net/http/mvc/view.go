package mvc

import (
	"context"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/alexfalkowski/go-service/meta"
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

// View to render.
type View string

// Render the view.
func (v View) Render(ctx context.Context, res http.ResponseWriter, model any) {
	if err := views.template.ExecuteTemplate(res, string(v), model); err != nil {
		meta.WithAttribute(ctx, "mvcViewError", meta.Error(err))
		res.WriteHeader(status.Code(err))
	}
}
