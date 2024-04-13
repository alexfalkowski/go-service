package debug

import (
	"github.com/felixge/fgprof"
)

// RegisterFgprof for debug.
func RegisterFgprof(server *Server) {
	server.Mux.Handle("/debug/fgprof", fgprof.Handler())
}
