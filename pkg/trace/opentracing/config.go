package opentracing

import (
	"github.com/alexfalkowski/go-service/pkg/trace/opentracing/datadog"
	"github.com/alexfalkowski/go-service/pkg/trace/opentracing/jaeger"
)

// Config for opentracing.
type Config struct {
	Datadog datadog.Config `yaml:"datadog"`
	Jaeger  jaeger.Config  `yaml:"jaeger"`
}
