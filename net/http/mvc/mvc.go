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

// RegisterParams defines dependencies used to register MVC package globals.
//
// This package uses package-level registration to avoid threading commonly shared dependencies
// through every routing helper call. These dependencies are typically provided by DI wiring.
type RegisterParams struct {
	di.In

	// Mux is the HTTP mux where MVC routes (Route/Get/Post/etc.) will be registered.
	Mux *http.ServeMux

	// FunctionMap is the template function map used when parsing templates.
	//
	// It is typically constructed via go-sprout and used to provide common template helpers.
	FunctionMap sprout.FunctionMap

	// FileSystem is the filesystem used to load template files and static files.
	//
	// It is optional: when nil, MVC is considered not defined and routing helpers return false.
	FileSystem fs.FS `optional:"true"`

	// Layout defines the base templates used to render full and partial views.
	//
	// It is optional: when nil or invalid (see Layout.IsValid), MVC is considered not defined and
	// routing helpers return false.
	Layout *Layout `optional:"true"`
}

// Register stores MVC dependencies in package-level variables.
//
// Register is expected to be called during application startup (typically via dependency injection).
//
// Definition rules:
// MVC routing/rendering is considered "defined" only when both:
//   - FileSystem is non-nil, and
//   - Layout is valid (see (*Layout).IsValid).
//
// If MVC is not defined, routing helpers (Route/Get/Post/etc. and static helpers) return false and do not
// register handlers.
func Register(params RegisterParams) {
	mux = params.Mux
	fmap = params.FunctionMap
	fileSystem = params.FileSystem
	layout = params.Layout
}

// IsDefined reports whether MVC routing and rendering has been configured.
//
// MVC is considered defined only when a FileSystem is available and Layout is valid (non-nil with both
// layout template names set).
func IsDefined() bool {
	return fileSystem != nil && layout.IsValid()
}
