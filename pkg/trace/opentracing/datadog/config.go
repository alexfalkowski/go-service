package datadog

import (
	"github.com/kelseyhightower/envconfig"
)

// NewConfig for datadog.
func NewConfig() (*Config, error) {
	var config Config
	err := envconfig.Process("", &config)

	return &config, err
}

// Config for datadog.
type Config struct {
	AppName   string `envconfig:"APP_NAME" required:"true"`
	TraceHost string `envconfig:"DATADOG_TRACE_HOST" required:"true" default:"localhost:8126"`
}
