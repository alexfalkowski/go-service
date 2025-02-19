package test

import (
	"github.com/alexfalkowski/go-service/runtime"
	"github.com/alexfalkowski/go-service/telemetry/logger"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"go.uber.org/fx"
)

// NewTracer for test.
func NewTracer(lc fx.Lifecycle, config *tracer.Config, logger *logger.Logger) *tracer.Tracer {
	params := tracer.Params{
		Lifecycle:   lc,
		Environment: Environment,
		Name:        Name,
		Version:     Version,
		FileSystem:  FS,
		Config:      config,
		Logger:      logger,
	}

	tracer, err := tracer.NewTracer(params)
	runtime.Must(err)

	return tracer
}
