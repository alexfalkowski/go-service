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

// NewLayout constructs a Layout that defines the base templates used for rendering views.
//
// The full parameter is the template file used as the base layout for "full page" renders.
// The partial parameter is the template file used as the base layout for "partial" renders.
//
// These values are later used by NewFullView/NewPartialView to parse templates from the
// registered filesystem (see mvc.Register and mvc.IsDefined).
func NewLayout(full, partial string) *Layout {
	return &Layout{full: full, partial: partial}
}

// Layout defines the base templates used to render full and partial views.
type Layout struct {
	full    string
	partial string
}

// Full returns the configured full layout template name/path.
func (l *Layout) Full() string {
	return l.full
}

// Partial returns the configured partial layout template name/path.
func (l *Layout) Partial() string {
	return l.partial
}

// FullName returns the base file name of the full layout template.
func (l *Layout) FullName() string {
	return l.name(l.full)
}

// PartialName returns the base file name of the partial layout template.
func (l *Layout) PartialName() string {
	return l.name(l.partial)
}

// IsValid reports whether l is non-nil and both layout template names are set.
//
// MVC is considered "defined" only when a filesystem is registered and the layout is valid.
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
//
// This is a convenience helper when a controller supports both full-page and partial rendering.
func NewViewPair(name string) (*View, *View) {
	return NewFullView(name), NewPartialView(name)
}

// NewFullView parses the full layout template and the view template from the registered filesystem.
//
// It uses the package-level filesystem, layout, and template function map registered via mvc.Register.
// If MVC is not defined, or if template parsing fails (missing files, parse errors), this function will
// panic because it uses template.Must.
func NewFullView(name string) *View {
	template := template.Must(template.New(strings.Empty).Funcs(fmap).ParseFS(fileSystem, layout.Full(), name))

	return &View{name: layout.FullName(), template: template}
}

// NewPartialView parses the partial layout template and the view template from the registered filesystem.
//
// It uses the package-level filesystem, layout, and template function map registered via mvc.Register.
// If MVC is not defined, or if template parsing fails (missing files, parse errors), this function will
// panic because it uses template.Must.
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
// Context requirements:
// Render expects the HTTP response writer to be present in ctx via net/http/meta.WithResponse.
// (Handlers created by this package's routing helpers populate that value before invoking controllers/views.)
//
// Render model:
// Render wraps the provided model in a Template which includes exported meta attributes under Template.Meta.
// This allows templates to access request-scoped metadata (for example requestId) without controllers having
// to explicitly thread those values through the model.
//
// Error handling:
// If template execution fails, Render records "mvcViewError" in ctx and writes a status code derived from
// the error. It does not write an error body.
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
