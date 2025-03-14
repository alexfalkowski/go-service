package debug

import "github.com/felixge/fgprof"

// RegisterFgprof for debug.
func RegisterFgprof(mux *ServeMux) {
	mux.Handle("/debug/fgprof", fgprof.Handler())
}
