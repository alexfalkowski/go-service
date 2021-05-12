package config

import (
	"github.com/kelseyhightower/envconfig"
)

// NewConfig for the services.
func NewConfig() (*Config, error) {
	var config Config
	err := envconfig.Process("", &config)

	return &config, err
}

// Config for the services.
type Config struct {
	AppName  string `envconfig:"APP_NAME" required:"true"`
	GRPCPort string `envconfig:"GRPC_PORT" required:"true" default:"8081"`
}
