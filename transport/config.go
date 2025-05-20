package transport

import (
	"github.com/alexfalkowski/go-service/v2/transport/grpc"
	"github.com/alexfalkowski/go-service/v2/transport/http"
)

// IsEnabled for transport.
func IsEnabled(cfg *Config) bool {
	return cfg != nil
}

// Config for transport.
type Config struct {
	GRPC *grpc.Config `yaml:"grpc,omitempty" json:"grpc,omitempty" toml:"grpc,omitempty"`
	HTTP *http.Config `yaml:"http,omitempty" json:"http,omitempty" toml:"http,omitempty"`
}
