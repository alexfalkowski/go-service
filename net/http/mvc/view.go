package mvc

import (
	"html/template"
	"io"
	"io/fs"

	"github.com/alexfalkowski/go-service/runtime"
)

// ParseTemplate for mvc.
func ParseTemplate(fs fs.FS, pattern string) *template.Template {
	t, err := template.ParseFS(fs, pattern)
	runtime.Must(err)

	return t
}

// NewView for routes in mvc.
func NewView(success, failure *template.Template) *View {
	return &View{success: success, failure: failure}
}

// View used for routes.
type View struct {
	success *template.Template
	failure *template.Template
}

// ExecuteSuccess for template.
func (t *View) ExecuteSuccess(wr io.Writer, data any) error {
	return execute(t.success, wr, data)
}

// ExecuteFailure for template.
func (t *View) ExecuteFailure(wr io.Writer, data any) error {
	return execute(t.failure, wr, data)
}

func execute(template *template.Template, wr io.Writer, data any) error {
	if template == nil {
		return nil
	}

	return template.Execute(wr, data)
}
