package debug

import (
	"net/http/pprof"

	"github.com/alexfalkowski/go-service/env"
	"github.com/arl/statsviz"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// RegisterParams for debug.
type RegisterParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    *Config
	Env       env.Environment
	Logger    *zap.Logger
}

// Register debug.
func Register(params RegisterParams) {
	if !params.Env.IsDevelopment() {
		return
	}

	m := mux(params.Lifecycle, params.Config, params.Logger)

	// Register statsviz.
	statsviz.Register(m)

	// Register pprof.
	m.HandleFunc("/debug/pprof/", pprof.Index)
	m.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	m.HandleFunc("/debug/pprof/profile", pprof.Profile)
	m.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	m.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// Register psutil.
	m.HandleFunc("/debug/psutil", psutil)
}
