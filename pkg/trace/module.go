package trace

import (
	"github.com/alexfalkowski/go-service/pkg/trace/opentracing/jaeger"
	"go.uber.org/fx"
)

var (
	// JaegerOpenTracing for fx.
	JaegerOpenTracing = fx.Invoke(jaeger.Register)
)
