package fgprof

import (
	"github.com/alexfalkowski/go-service/v2/debug/http"
	"github.com/felixge/fgprof"
)

// Register for debug.
func Register(mux *http.ServeMux) {
	mux.Handle("/debug/fgprof", fgprof.Handler())
}
