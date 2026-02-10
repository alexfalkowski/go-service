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

// RegisterParams defines dependencies used to register MVC globals.
type RegisterParams struct {
	di.In
	Mux         *http.ServeMux
	FunctionMap sprout.FunctionMap
	FileSystem  fs.FS   `optional:"true"`
	Layout      *Layout `optional:"true"`
}

// Register stores the MVC dependencies used by the routing and view helpers.
//
// MVC routes are considered defined only when both a FileSystem and a valid Layout have been registered.
// It is expected that Register is called during application startup (typically via Fx).
func Register(params RegisterParams) {
	mux = params.Mux
	fmap = params.FunctionMap
	fileSystem = params.FileSystem
	layout = params.Layout
}

// IsDefined reports whether MVC routing and rendering has been configured.
//
// MVC is considered defined only when a FileSystem is available and Layout is valid.
func IsDefined() bool {
	return fileSystem != nil && layout.IsValid()
}
