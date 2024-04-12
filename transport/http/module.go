package http

import (
	"github.com/alexfalkowski/go-service/transport/http/telemetry/metrics"
	"github.com/alexfalkowski/go-service/transport/http/telemetry/tracer"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	fx.Provide(NewServeMux),
	fx.Provide(NewServer),
	fx.Provide(tracer.NewTracer),
	fx.Invoke(metrics.Register),
	fx.Provide(metrics.NewMeter),
)
