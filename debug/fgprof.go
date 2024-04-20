package debug

import (
	"net/http"

	"github.com/felixge/fgprof"
)

// RegisterFgprof for debug.
func RegisterFgprof(mux *http.ServeMux) {
	mux.Handle("/debug/fgprof", fgprof.Handler())
}
