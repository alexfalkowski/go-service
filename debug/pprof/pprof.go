package pprof

import (
	"net/http/pprof"

	"github.com/alexfalkowski/go-service/v2/debug/http"
	"github.com/alexfalkowski/go-service/v2/env"
)

// Register installs net/http/pprof handlers on mux.
//
// Handlers are registered under the "/debug/pprof" prefix (namespaced by service name via http.Pattern),
// providing standard pprof endpoints such as:
//   - index and profile listing
//   - cmdline
//   - CPU profile
//   - symbol lookup
//   - trace
//
// This registration is intended to be composed into the go-service debug server wiring.
func Register(name env.Name, mux *http.ServeMux) {
	mux.HandleFunc(http.Pattern(name, "/debug/pprof/"), pprof.Index)
	mux.HandleFunc(http.Pattern(name, "/debug/pprof/cmdline"), pprof.Cmdline)
	mux.HandleFunc(http.Pattern(name, "/debug/pprof/profile"), pprof.Profile)
	mux.HandleFunc(http.Pattern(name, "/debug/pprof/symbol"), pprof.Symbol)
	mux.HandleFunc(http.Pattern(name, "/debug/pprof/trace"), pprof.Trace)
}
