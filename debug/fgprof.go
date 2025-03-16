package debug

import "github.com/felixge/fgprof"

// RegisterFgprof for debug.
func RegisterFgprof(srv *Server) {
	if srv == nil {
		return
	}

	mux := srv.ServeMux()

	mux.Handle("/debug/fgprof", fgprof.Handler())
}
