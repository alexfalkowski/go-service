package debug

import (
	"net/http/pprof"

	"github.com/alexfalkowski/go-service/env"
	"github.com/arl/statsviz"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Register debug.
func Register(lc fx.Lifecycle, env env.Environment, logger *zap.Logger) {
	if !env.IsDevelopment() {
		return
	}

	mux := newServeMux(lc, logger)

	// Register statsviz.
	statsviz.Register(mux)

	// Register pprof.
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// Register psutil.
	mux.HandleFunc("/debug/psutil", psutil)
}
