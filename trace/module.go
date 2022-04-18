package trace

import (
	"github.com/alexfalkowski/go-service/trace/opentracing"
	"go.uber.org/fx"
)

var (
	// JaegerOpenTracingModule for fx.
	// nolint:gochecknoglobals
	JaegerOpenTracingModule = fx.Options(
		fx.Provide(opentracing.NewJaegerDatabaseTracer),
		fx.Provide(opentracing.NewJaegerCacheTracer),
		fx.Provide(opentracing.NewJaegerTransportTracer),
	)

	// DataDogOpenTracingModule for fx.
	// nolint:gochecknoglobals
	DataDogOpenTracingModule = fx.Options(
		fx.Provide(opentracing.NewDataDogDatabaseTracer),
		fx.Provide(opentracing.NewDataDogCacheTracer),
		fx.Provide(opentracing.NewDataDogTransportTracer),
	)
)
