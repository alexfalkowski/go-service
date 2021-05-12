package trace

import (
	"github.com/alexfalkowski/go-service/pkg/trace/opentracing/datadog"
	"github.com/alexfalkowski/go-service/pkg/trace/opentracing/jaeger"
	"go.uber.org/fx"
)

var (
	// JaegerOpenTracing for fx.
	JaegerOpenTracing = fx.Options(fx.Invoke(jaeger.Register), fx.Provide(jaeger.NewConfig))

	// DataDogOpenTracing for fx.
	DataDogOpenTracing = fx.Options(fx.Invoke(datadog.Register), fx.Provide(datadog.NewConfig))
)
