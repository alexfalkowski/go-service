package trace

import (
	"github.com/alexfalkowski/go-service/pkg/trace/opentracing"
	"go.uber.org/fx"
)

var (
	// JaegerOpenTracing for fx.
	JaegerOpenTracing = fx.Invoke(opentracing.RegisterJaeger)

	// Module for fx.
	Module = fx.Options(JaegerOpenTracing)
)
