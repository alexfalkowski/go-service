package test

import (
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"go.uber.org/fx"
)

// NewTracer for test.
func NewTracer(lc fx.Lifecycle, config *tracer.Config) *tracer.Tracer {
	params := tracer.TracerParams{
		Lifecycle:   lc,
		Environment: Environment,
		Name:        Name,
		Version:     Version,
		Config:      config,
	}

	return tracer.NewTracer(params)
}
