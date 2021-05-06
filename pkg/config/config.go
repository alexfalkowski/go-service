package config

import (
	"github.com/kelseyhightower/envconfig"
)

// NewConfig for the services.
func NewConfig() (*Config, error) {
	var config Config

	if err := envconfig.Process("", &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// Config for the services.
type Config struct {
	AppName         string `envconfig:"APP_NAME" required:"true"`
	HTTPPort        string `envconfig:"HTTP_PORT" required:"true" default:"8080"`
	GRPCPort        string `envconfig:"GRPC_PORT" required:"true" default:"8081"`
	DatabaseURL     string `envconfig:"DATABASE_URL" required:"true"`
	JaegerTraceHost string `envconfig:"JAEGER_TRACE_HOST" required:"true" default:"localhost:6831"`
	CacheHost       string `envconfig:"CACHE_HOST" required:"true" default:"localhost:6379"`
}
