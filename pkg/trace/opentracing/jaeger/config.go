package jaeger

import (
	"github.com/kelseyhightower/envconfig"
)

// NewConfig for jaeger.
func NewConfig() (*Config, error) {
	var config Config
	err := envconfig.Process("", &config)

	return &config, err
}

// Config for jaeger.
type Config struct {
	AppName   string `envconfig:"APP_NAME" required:"true"`
	TraceHost string `envconfig:"JAEGER_TRACE_HOST" required:"true" default:"localhost:6831"`
}
