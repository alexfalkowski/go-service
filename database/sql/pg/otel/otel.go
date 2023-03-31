package otel

import (
	"github.com/alexfalkowski/go-service/otel"
	"github.com/alexfalkowski/go-service/version"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

// TracerParams for otel.
type TracerParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    *otel.Config
	Version   version.Version
}

// NewTracer for otel.
func NewTracer(params TracerParams) (Tracer, error) {
	return otel.NewTracer(otel.TracerParams{Lifecycle: params.Lifecycle, Name: "pg", Config: params.Config, Version: params.Version})
}

// Tracer for otel.
type Tracer trace.Tracer
