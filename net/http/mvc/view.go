package mvc

import (
	"html/template"
	"io/fs"
	"path"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/net/http"
	httpmeta "github.com/alexfalkowski/go-service/v2/net/http/meta"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-sync"
)

// ErrMissingView is returned when MVC rendering is requested without a view.
var ErrMissingView = errors.New("mvc: missing view")

// NewLayout constructs a Layout that defines the base templates used for rendering views.
//
// The full parameter is the template file used as the base layout for "full page" renders.
// The partial parameter is the template file used as the base layout for "partial" renders.
//
// These values are later used to parse templates from the configured filesystem.
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
	return path.Base(name)
}

// NewViewPair returns a full and partial View pair for name.
//
// This is a convenience helper when a controller supports both full-page and partial rendering.
// Call it during startup or route registration so template read and parse failures fail fast.
func (s *Server) NewViewPair(name string) (*View, *View) {
	return s.NewFullView(name), s.NewPartialView(name)
}

// NewFullView parses the full layout template and the view template from the configured filesystem.
//
// Call it during startup or route registration so template read and parse failures fail fast.
// If MVC is not defined, or if template parsing fails (missing files, parse errors), this method will
// panic.
func (s *Server) NewFullView(name string) *View {
	return s.newView(s.layout.Full(), name)
}

// NewPartialView parses the partial layout template and the view template from the configured filesystem.
//
// Call it during startup or route registration so template read and parse failures fail fast.
// If MVC is not defined, or if template parsing fails (missing files, parse errors), this method will
// panic.
func (s *Server) NewPartialView(name string) *View {
	return s.newView(s.layout.Partial(), name)
}

func (s *Server) newView(layoutName, name string) *View {
	tmpl := template.New(layoutName).Funcs(s.fmap)
	s.parseTemplate(tmpl, layoutName)
	s.parseTemplate(tmpl.New(name), name)

	return &View{name: layoutName, template: tmpl, pool: s.pool}
}

func (s *Server) parseTemplate(tmpl *template.Template, name string) {
	data, err := fs.ReadFile(s.fileSystem, name)
	runtime.Must(err)

	_, err = tmpl.Parse(string(data))
	runtime.Must(err)
}

// View renders an HTML template.
type View struct {
	template *template.Template
	pool     *sync.BufferPool
	name     string
}

// Render executes the view template against a Template model and writes it to the HTTP response writer.
//
// Context requirements:
// Render expects the HTTP response writer to be present in ctx via net/http/meta.
// (Handlers created by this package's routing helpers populate that value before invoking controllers/views.)
//
// Render model:
// Render wraps the provided model in a Template which includes exported meta attributes under [Template.Meta].
// This allows templates to access request-scoped metadata (for example requestId) without controllers having
// to explicitly thread those values through the model.
//
// Error handling:
// Render renders into an internal buffer before writing to the response so template execution failures do not
// commit a partial response. If ctx is already canceled, Render returns the context error without executing the
// template. Otherwise, Render returns any template execution error or final response write error.
//
// Size limits:
// MVC does not enforce body size caps itself. In supported wiring, inbound request bodies are capped by the
// transport HTTP server before MVC handlers run, and outbound response bodies are capped by go-service HTTP
// clients when responses are read.
func (v *View) Render(ctx context.Context, model any) error {
	buffer := v.pool.Get()
	defer v.pool.Put(buffer)

	if err := v.render(ctx, buffer, model); err != nil {
		return err
	}

	_, err := buffer.WriteTo(httpmeta.Response(ctx))
	return err
}

func (v *View) render(ctx context.Context, writer io.Writer, model any) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	template := &Template{
		Meta:  meta.Strings(ctx, meta.NoPrefix),
		Model: model,
	}

	return v.template.ExecuteTemplate(writer, v.name, template)
}

func (s *Server) writeView(ctx context.Context, res http.ResponseWriter, view *View, model any, code int) {
	if err := s.renderView(ctx, res, view, model, code); err != nil {
		res.WriteHeader(status.Code(err))
	}
}

func (s *Server) renderView(ctx context.Context, res http.ResponseWriter, view *View, model any, code int) error {
	if view == nil {
		status.RecordError(ctx, ErrMissingView)
		return ErrMissingView
	}

	buffer := s.pool.Get()
	defer s.pool.Put(buffer)

	if err := view.render(ctx, buffer, model); err != nil {
		status.RecordError(ctx, err)
		return err
	}

	writeBuffer(res, code, buffer)
	return nil
}

func writeBuffer(res http.ResponseWriter, code int, buffer *bytes.Buffer) {
	res.WriteHeader(code)
	_, _ = buffer.WriteTo(res)
}
