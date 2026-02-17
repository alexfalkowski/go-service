package pprof

import (
	"net/http/pprof"

	"github.com/alexfalkowski/go-service/v2/debug/http"
	"github.com/alexfalkowski/go-service/v2/env"
)

// Register installs pprof handlers under /debug/pprof.
func Register(name env.Name, mux *http.ServeMux) {
	mux.HandleFunc(http.Pattern(name, "/debug/pprof/"), pprof.Index)
	mux.HandleFunc(http.Pattern(name, "/debug/pprof/cmdline"), pprof.Cmdline)
	mux.HandleFunc(http.Pattern(name, "/debug/pprof/profile"), pprof.Profile)
	mux.HandleFunc(http.Pattern(name, "/debug/pprof/symbol"), pprof.Symbol)
	mux.HandleFunc(http.Pattern(name, "/debug/pprof/trace"), pprof.Trace)
}
