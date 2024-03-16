package debug

import (
	"net/http/pprof"

	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/marshaller"
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
	JSON      *marshaller.JSON
	Logger    *zap.Logger
}

// Register debug.
func Register(params RegisterParams) {
	if !params.Env.IsDevelopment() {
		return
	}

	s := newServer(params.Lifecycle, params.Config, params.JSON, params.Logger)

	// Register statsviz.
	statsviz.Register(s.mux)

	// Register pprof.
	s.mux.HandleFunc("/debug/pprof/", pprof.Index)
	s.mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	s.mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	s.mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	s.mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// Register psutil.
	s.mux.HandleFunc("/debug/psutil", s.psutil)
}
