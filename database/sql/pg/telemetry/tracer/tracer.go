package tracer

import (
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/version"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

// Params for tracer.
type Params struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    *tracer.Config
	Version   version.Version
}

// NewTracer for tracer.
func NewTracer(params Params) (Tracer, error) {
	return tracer.NewTracer(tracer.Params{Lifecycle: params.Lifecycle, Name: "pg", Config: params.Config, Version: params.Version})
}

// Tracer for tracer.
type Tracer trace.Tracer
