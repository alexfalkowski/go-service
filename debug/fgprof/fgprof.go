package fgprof

import (
	"github.com/alexfalkowski/go-service/v2/debug/http"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/felixge/fgprof"
)

// Register installs the fgprof handler under /debug/fgprof.
func Register(name env.Name, mux *http.ServeMux) {
	mux.Handle(http.Pattern(name, "/debug/fgprof"), fgprof.Handler())
}
