package grpc

import (
	"github.com/kelseyhightower/envconfig"
)

// NewConfig for gRPC.
func NewConfig() (*Config, error) {
	var config Config
	err := envconfig.Process("", &config)

	return &config, err
}

// Config for gRPC.
type Config struct {
	Port string `envconfig:"GRPC_PORT" required:"true" default:"8081"`
}
