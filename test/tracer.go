package test

import (
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// NewTracer for test.
func NewTracer(lc fx.Lifecycle, logger *zap.Logger) trace.Tracer {
	tracer, err := tracer.NewTracer(lc, Environment, Version, NewOTLPTracerConfig(), logger)
	if err != nil {
		panic(err)
	}

	return tracer
}
