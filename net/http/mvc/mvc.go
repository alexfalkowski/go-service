package mvc

import (
	"io/fs"

	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/go-sprout/sprout"
)

var (
	mux        *http.ServeMux
	fmap       sprout.FunctionMap
	fileSystem fs.FS
	layout     *Layout
)

// RegisterParams for mvc.
type RegisterParams struct {
	di.In

	Mux         *http.ServeMux
	FunctionMap sprout.FunctionMap
	FileSystem  fs.FS   `optional:"true"`
	Layout      *Layout `optional:"true"`
}

// Register for mvc.
func Register(params RegisterParams) {
	mux = params.Mux
	fmap = params.FunctionMap
	fileSystem = params.FileSystem
	layout = params.Layout
}

// IsDefined for mvc.
func IsDefined() bool {
	return fileSystem != nil && layout.IsValid()
}
