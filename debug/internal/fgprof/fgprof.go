package fgprof

import (
	"github.com/alexfalkowski/go-service/v2/debug/http"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/felixge/fgprof"
)

// Register installs the fgprof handler on mux.
//
// The handler is registered at "/debug/fgprof" (namespaced by service name via http.Pattern) and serves
// fgprof's wall-clock based profiling UI/data. This is useful for diagnosing CPU usage as well as time
// spent blocked or waiting in scheduling.
//
// This registration is intended to be composed into the go-service debug server wiring.
func Register(name env.Name, mux *http.ServeMux) {
	mux.Handle(http.Pattern(name, "/debug/fgprof"), fgprof.Handler())
}
