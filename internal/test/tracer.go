package test

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
)

// RegisterTracer installs the shared test tracer provider on the supplied lifecycle.
func RegisterTracer(lc di.Lifecycle, config *tracer.Config) {
	params := tracer.TracerParams{
		Lifecycle:   lc,
		Environment: Environment,
		Name:        Name,
		Version:     Version,
		Config:      config,
	}

	runtime.Must(tracer.Register(params))
}
