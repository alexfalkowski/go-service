package transport

import (
	"github.com/alexfalkowski/go-service/v2/transport/grpc"
	"github.com/alexfalkowski/go-service/v2/transport/http"
)

// Config for transport.
type Config struct {
	GRPC *grpc.Config `yaml:"grpc,omitempty" json:"grpc,omitempty" toml:"grpc,omitempty"`
	HTTP *http.Config `yaml:"http,omitempty" json:"http,omitempty" toml:"http,omitempty"`
}

// IsEnabled for transport.
func (c *Config) IsEnabled() bool {
	return c != nil
}
