package trace

import (
	"github.com/alexfalkowski/go-service/pkg/trace/opentracing"
)

// Config for trace.
type Config struct {
	Opentracing opentracing.Config `yaml:"opentracing"`
}
