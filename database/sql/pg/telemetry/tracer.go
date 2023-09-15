package telemetry

import (
	"github.com/alexfalkowski/go-service/telemetry"
	"github.com/alexfalkowski/go-service/version"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

// TracerParams for otel.
type TracerParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    *telemetry.Config
	Version   version.Version
}

// NewTracer for otel.
func NewTracer(params TracerParams) (Tracer, error) {
	return telemetry.NewTracer(telemetry.TracerParams{Lifecycle: params.Lifecycle, Name: "pg", Config: params.Config, Version: params.Version})
}

// Tracer for otel.
type Tracer trace.Tracer
