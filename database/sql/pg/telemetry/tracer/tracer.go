package tracer

import (
	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/version"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

// Params for tracer.
type Params struct {
	fx.In

	Lifecycle   fx.Lifecycle
	Config      *tracer.Config
	Environment env.Environment
	Version     version.Version
}

// NewTracer for tracer.
func NewTracer(params Params) (Tracer, error) {
	return tracer.NewTracer(params.Lifecycle, "pg", params.Environment, params.Version, params.Config)
}

// Tracer for tracer.
type Tracer trace.Tracer
