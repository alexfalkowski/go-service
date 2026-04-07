package test

import (
	"github.com/alexfalkowski/go-service/v2/debug"
	"github.com/alexfalkowski/go-service/v2/debug/fgprof"
	"github.com/alexfalkowski/go-service/v2/debug/http"
	"github.com/alexfalkowski/go-service/v2/debug/pprof"
	"github.com/alexfalkowski/go-service/v2/debug/psutil"
	"github.com/alexfalkowski/go-service/v2/debug/statsviz"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
)

// NewDebugServer returns a debug server with the standard pprof, fgprof, psutil, and statsviz handlers registered.
func NewDebugServer(config *debug.Config, logger *logger.Logger) *debug.Server {
	server, err := newDebugServer(config, logger)
	runtime.Must(err)

	return server
}

func newDebugServer(config *debug.Config, logger *logger.Logger) (*debug.Server, error) {
	mux := http.NewServeMux()
	server, err := debug.NewServer(debug.ServerParams{
		Shutdowner: NewShutdowner(),
		Mux:        mux,
		Config:     config,
		Logger:     logger,
		FS:         FS,
	})
	if err != nil {
		return nil, err
	}

	pprof.Register(Name, mux)
	fgprof.Register(Name, mux)
	psutil.Register(Name, Content, mux)

	err = statsviz.Register(Name, mux)
	if err != nil {
		return nil, err
	}

	return server, nil
}
