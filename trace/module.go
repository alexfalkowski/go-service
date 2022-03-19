package trace

import (
	"github.com/alexfalkowski/go-service/trace/opentracing/datadog"
	"github.com/alexfalkowski/go-service/trace/opentracing/jaeger"
	"go.uber.org/fx"
)

var (
	// JaegerOpenTracingModule for fx.
	// nolint:gochecknoglobals
	JaegerOpenTracingModule = fx.Options(fx.Invoke(jaeger.Register))

	// DataDogOpenTracingModule for fx.
	// nolint:gochecknoglobals
	DataDogOpenTracingModule = fx.Options(fx.Invoke(datadog.Register))
)
