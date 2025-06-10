package test

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
)

// NewTracer for test.
func NewTracer(lc di.Lifecycle, config *tracer.Config) *tracer.Tracer {
	params := tracer.TracerParams{
		Lifecycle:   lc,
		Environment: Environment,
		Name:        Name,
		Version:     Version,
		Config:      config,
	}

	return tracer.NewTracer(params)
}
