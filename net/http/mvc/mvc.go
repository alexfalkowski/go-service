package mvc

import (
	"io/fs"
	"net/http"

	"go.uber.org/fx"
)

var (
	mux        *http.ServeMux
	fileSystem fs.FS
	layout     *Layout
)

// RegisterParams for mvc.
type RegisterParams struct {
	fx.In

	Mux        *http.ServeMux
	FileSystem fs.FS   `optional:"true"`
	Layout     *Layout `optional:"true"`
}

// Register for mvc.
func Register(params RegisterParams) {
	mux, fileSystem, layout = params.Mux, params.FileSystem, params.Layout
}

// IsDefined for mvc.
func IsDefined() bool {
	return fileSystem != nil && layout != nil && layout.IsValid()
}
