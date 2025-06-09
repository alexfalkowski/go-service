package pprof

import (
	"net/http/pprof"

	"github.com/alexfalkowski/go-service/v2/debug/http"
	"github.com/alexfalkowski/go-service/v2/env"
)

// Register for debug.
func Register(name env.Name, mux *http.ServeMux) {
	mux.HandleFunc(http.Pattern("/debug/pprof/", name), pprof.Index)
	mux.HandleFunc(http.Pattern("/debug/pprof/cmdline", name), pprof.Cmdline)
	mux.HandleFunc(http.Pattern("/debug/pprof/profile", name), pprof.Profile)
	mux.HandleFunc(http.Pattern("/debug/pprof/symbol", name), pprof.Symbol)
	mux.HandleFunc(http.Pattern("/debug/pprof/trace", name), pprof.Trace)
}
