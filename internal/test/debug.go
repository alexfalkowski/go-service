package test

import (
	"github.com/alexfalkowski/go-service/v2/debug"
	"github.com/alexfalkowski/go-service/v2/debug/http"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
)

// NewDebugServer returns a debug server with the standard debug handlers registered.
func NewDebugServer(lc di.Lifecycle, config *debug.Config, logger *logger.Logger) (*debug.Server, error) {
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

	if err := debug.Register(debug.RegisterParams{
		Config:    config,
		Lifecycle: lc,
		Name:      Name,
		Content:   Content,
		Mux:       mux,
	}); err != nil {
		return nil, err
	}

	return server, nil
}
