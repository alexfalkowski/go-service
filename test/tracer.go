package test

import (
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

// NewTracer for test.
func NewTracer(lc fx.Lifecycle) trace.Tracer {
	tracer, err := tracer.NewTracer(lc, Environment, Version, NewOTLPTracerConfig())
	if err != nil {
		panic(err)
	}

	return tracer
}
