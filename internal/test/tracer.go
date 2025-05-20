package test

import (
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"go.uber.org/fx"
)

// NewTracer for test.
func NewTracer(lc fx.Lifecycle, config *tracer.Config) *tracer.Tracer {
	params := tracer.Params{
		Lifecycle:   lc,
		Environment: Environment,
		Name:        Name,
		Version:     Version,
		FileSystem:  FS,
		Config:      config,
	}

	tracer, err := tracer.NewTracer(params)
	runtime.Must(err)

	return tracer
}
