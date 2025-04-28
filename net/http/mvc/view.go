package mvc

import (
	"context"
	"html/template"
	"path/filepath"

	"github.com/alexfalkowski/go-service/meta"
	hm "github.com/alexfalkowski/go-service/net/http/meta"
	"github.com/alexfalkowski/go-service/net/http/status"
	"github.com/alexfalkowski/go-service/strings"
	"github.com/go-sprout/sprout/sprigin"
)

// NewLayout defines a full and partial layout.
func NewLayout(full, partial string) *Layout {
	return &Layout{full: full, partial: partial}
}

// Layout is the main template that is used for all other templates.
type Layout struct {
	full    string
	partial string
}

// Full of the layout.
func (l *Layout) Full() string {
	return l.full
}

// Partial of the layout.
func (l *Layout) Partial() string {
	return l.partial
}

// FullName of the layout.
func (l *Layout) FullName() string {
	return l.name(l.full)
}

// PartialName of the layout.
func (l *Layout) PartialName() string {
	return l.name(l.partial)
}

// IsValid the layout.
func (l *Layout) IsValid() bool {
	return !strings.IsEmpty(l.full) || !strings.IsEmpty(l.partial)
}

func (l *Layout) name(name string) string {
	return filepath.Base(name)
}

// NewView to render.
func NewView(name string) *View {
	template := template.Must(template.New("").Funcs(sprigin.FuncMap()).ParseFS(fileSystem, layout.Full(), name))

	return &View{name: layout.FullName(), template: template}
}

// NewPartialView to render.
func NewPartialView(name string) *View {
	template := template.Must(template.New("").Funcs(sprigin.FuncMap()).ParseFS(fileSystem, layout.Partial(), name))

	return &View{name: layout.PartialName(), template: template}
}

// View to render.
type View struct {
	template *template.Template
	name     string
}

// Render the view.
func (v *View) Render(ctx context.Context, model any) {
	res := hm.Response(ctx)

	if err := v.template.ExecuteTemplate(res, v.name, model); err != nil {
		meta.WithAttribute(ctx, "mvcViewError", meta.Error(err))
		res.WriteHeader(status.Code(err))
	}
}
