package test

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
)

// RegisterTracer for test.
func RegisterTracer(lc di.Lifecycle, config *tracer.Config) {
	params := tracer.TracerParams{
		Lifecycle:   lc,
		Environment: Environment,
		Name:        Name,
		Version:     Version,
		Config:      config,
	}

	tracer.Register(params)
}
