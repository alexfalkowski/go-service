package mvc

import (
	"context"
	"html/template"
	"path/filepath"

	"github.com/alexfalkowski/go-service/meta"
	hm "github.com/alexfalkowski/go-service/net/http/meta"
	"github.com/alexfalkowski/go-service/net/http/status"
	"github.com/go-sprout/sprout/sprigin"
)

// Layout is the main template that is used for all other templates.
type Layout string

// String of the layout.
func (l Layout) String() string {
	return string(l)
}

// Name of the layout.
func (l Layout) Name() string {
	return filepath.Base(l.String())
}

// NewView to render.
func NewView(name string) *View {
	template := template.Must(template.New("").Funcs(sprigin.FuncMap()).ParseFS(fileSystem, layout.String(), name))

	return &View{template: template}
}

// View to render.
type View struct {
	template *template.Template
}

// Render the view.
func (v *View) Render(ctx context.Context, model any) {
	res := hm.Response(ctx)

	if err := v.template.ExecuteTemplate(res, layout.Name(), model); err != nil {
		meta.WithAttribute(ctx, "mvcViewError", meta.Error(err))
		res.WriteHeader(status.Code(err))
	}
}
