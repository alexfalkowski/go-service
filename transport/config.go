package transport

import (
	"github.com/alexfalkowski/go-service/structs"
	"github.com/alexfalkowski/go-service/transport/grpc"
	"github.com/alexfalkowski/go-service/transport/http"
)

// IsEnabled for transport.
func IsEnabled(cfg *Config) bool {
	return !structs.IsZero(cfg)
}

// Config for transport.
type Config struct {
	GRPC *grpc.Config `yaml:"grpc,omitempty" json:"grpc,omitempty" toml:"grpc,omitempty"`
	HTTP *http.Config `yaml:"http,omitempty" json:"http,omitempty" toml:"http,omitempty"`
}
