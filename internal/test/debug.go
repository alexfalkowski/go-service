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

// NewDebugServer for test.
func NewDebugServer(config *debug.Config, logger *logger.Logger) *debug.Server {
	mux := http.NewServeMux()
	server, err := debug.NewServer(debug.ServerParams{
		Shutdowner: NewShutdowner(),
		Mux:        mux,
		Config:     config,
		Logger:     logger,
		FS:         FS,
	})
	runtime.Must(err)

	pprof.Register(mux)
	fgprof.Register(mux)
	psutil.Register(mux, Content)

	err = statsviz.Register(mux)
	runtime.Must(err)

	return server
}
