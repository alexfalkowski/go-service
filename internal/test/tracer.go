package test

import (
	"github.com/alexfalkowski/go-service/runtime"
	"github.com/alexfalkowski/go-service/telemetry/logger"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

// NewTracer for test.
func NewTracer(lc fx.Lifecycle, config *tracer.Config, logger *logger.Logger) trace.Tracer {
	tracer, err := tracer.NewTracer(lc, Environment, Version, Name, FS, config, logger)
	runtime.Must(err)

	return tracer
}
