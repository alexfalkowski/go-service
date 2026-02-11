package transport

import (
	"github.com/alexfalkowski/go-service/v2/transport/grpc"
	"github.com/alexfalkowski/go-service/v2/transport/http"
)

// Config configures service transports.
type Config struct {
	// GRPC configures gRPC transport behavior (servers/clients), if enabled.
	GRPC *grpc.Config `yaml:"grpc,omitempty" json:"grpc,omitempty" toml:"grpc,omitempty"`

	// HTTP configures HTTP transport behavior (servers/clients), if enabled.
	HTTP *http.Config `yaml:"http,omitempty" json:"http,omitempty" toml:"http,omitempty"`
}

// IsEnabled for transport.
func (c *Config) IsEnabled() bool {
	return c != nil
}
