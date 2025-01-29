package grpc

import (
	"github.com/alexfalkowski/go-service/server"
	"github.com/alexfalkowski/go-service/structs"
)

// IsEnabled for gRPC.
func IsEnabled(cfg *Config) bool {
	return !structs.IsZero(cfg)
}

// Config for gRPC.
type Config struct {
	*server.Config `yaml:",inline" json:",inline" toml:",inline"`
}
