package trace

import (
	"github.com/alexfalkowski/go-service/pkg/trace/opentracing/datadog"
	"github.com/alexfalkowski/go-service/pkg/trace/opentracing/jaeger"
	"go.uber.org/fx"
)

var (
	// JaegerOpenTracingModule for fx.
	JaegerOpenTracingModule = fx.Options(fx.Invoke(jaeger.Register))

	// DataDogOpenTracingModule for fx.
	DataDogOpenTracingModule = fx.Options(fx.Invoke(datadog.Register))
)
