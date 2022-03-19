package opentracing

import (
	"github.com/alexfalkowski/go-service/trace/opentracing/datadog"
	"github.com/alexfalkowski/go-service/trace/opentracing/jaeger"
)

// Config for opentracing.
type Config struct {
	Datadog datadog.Config `yaml:"datadog"`
	Jaeger  jaeger.Config  `yaml:"jaeger"`
}
