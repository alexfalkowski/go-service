package debug

import (
	"net/http/pprof"
)

// RegisterPprof for debug.
func RegisterPprof(server *Server) {
	server.Mux.HandleFunc("/debug/pprof/", pprof.Index)
	server.Mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	server.Mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	server.Mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	server.Mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
}
