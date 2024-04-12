package grpc

import (
	"github.com/alexfalkowski/go-service/transport/grpc/telemetry/metrics"
	"github.com/alexfalkowski/go-service/transport/grpc/telemetry/tracer"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	fx.Provide(NewServer),
	fx.Provide(tracer.NewTracer),
	fx.Provide(metrics.NewMeter),
)
