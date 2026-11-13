package debug

import (
	"github.com/alexfalkowski/go-service/v3/debug/http"
	"github.com/alexfalkowski/go-service/v3/debug/internal/fgprof"
	"github.com/alexfalkowski/go-service/v3/debug/internal/pprof"
	"github.com/alexfalkowski/go-service/v3/debug/internal/statsviz"
	"github.com/alexfalkowski/go-service/v3/env"
	"github.com/alexfalkowski/go-signal/v2"
)

// RegisterParams defines dependencies for registering built-in debug endpoints.
//
// It is a go-signal parameter struct ([signal.In]) used to install the built-in debug handlers
// on the debug mux when debug is enabled.
type RegisterParams struct {
	signal.In

	// Lifecycle owns endpoint resources that need shutdown, such as statsviz.
	Lifecycle signal.Lifecycle

	// Config enables or disables debug endpoint registration.
	Config *Config

	// Mux is the debug mux where handlers are registered.
	Mux *http.ServeMux

	// Name is the service name used to prefix debug routes.
	Name env.Name
}

// Register installs the built-in debug endpoint handlers when cfg is enabled.
//
// This is the public front door for debug endpoint registration. It registers
// pprof, fgprof, and statsviz handlers on Mux. When Config is nil or
// disabled, Register returns without installing handlers or starting any
// endpoint-owned background work.
func Register(params RegisterParams) error {
	if !params.Config.IsEnabled() {
		return nil
	}

	pprof.Register(params.Name, params.Mux)
	fgprof.Register(params.Name, params.Mux)

	return statsviz.Register(params.Lifecycle, params.Name, params.Mux)
}
