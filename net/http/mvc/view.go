package mvc

import (
	"html/template"
	"path/filepath"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/meta"
	hm "github.com/alexfalkowski/go-service/v2/net/http/meta"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// NewLayout constructs a Layout defining the full and partial layout template names.
func NewLayout(full, partial string) *Layout {
	return &Layout{full: full, partial: partial}
}

// Layout defines the base templates used to render full and partial views.
type Layout struct {
	full    string
	partial string
}

// Full returns the configured full layout template name.
func (l *Layout) Full() string {
	return l.full
}

// Partial returns the configured partial layout template name.
func (l *Layout) Partial() string {
	return l.partial
}

// FullName returns the base name of the full layout template.
func (l *Layout) FullName() string {
	return l.name(l.full)
}

// PartialName returns the base name of the partial layout template.
func (l *Layout) PartialName() string {
	return l.name(l.partial)
}

// IsValid reports whether l is non-nil and both layout template names are set.
func (l *Layout) IsValid() bool {
	if l == nil {
		return false
	}

	return !strings.IsEmpty(l.full) && !strings.IsEmpty(l.partial)
}

func (l *Layout) name(name string) string {
	return filepath.Base(name)
}

// NewViewPair returns a full and partial View pair for name.
func NewViewPair(name string) (*View, *View) {
	return NewFullView(name), NewPartialView(name)
}

// NewFullView parses the full layout and name templates from the registered filesystem.
func NewFullView(name string) *View {
	template := template.Must(template.New(strings.Empty).Funcs(fmap).ParseFS(fileSystem, layout.Full(), name))

	return &View{name: layout.FullName(), template: template}
}

// NewPartialView parses the partial layout and name templates from the registered filesystem.
func NewPartialView(name string) *View {
	template := template.Must(template.New(strings.Empty).Funcs(fmap).ParseFS(fileSystem, layout.Partial(), name))

	return &View{name: layout.PartialName(), template: template}
}

// View renders an HTML template.
type View struct {
	template *template.Template
	name     string
}

// Render executes the view template against a Template model and writes it to the HTTP response writer.
//
// Render expects the HTTP response writer to be present in ctx via net/http/meta.WithResponse.
// If template execution fails, it records "mvcViewError" in ctx and writes a status code derived from the error.
func (v *View) Render(ctx context.Context, model any) {
	res := hm.Response(ctx)
	template := &Template{
		Meta:  meta.Strings(ctx, meta.NoPrefix),
		Model: model,
	}

	if err := v.template.ExecuteTemplate(res, v.name, template); err != nil {
		meta.WithAttribute(ctx, "mvcViewError", meta.Error(err))
		res.WriteHeader(status.Code(err))
	}
}
