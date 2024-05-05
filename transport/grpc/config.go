package grpc

import (
	"github.com/alexfalkowski/go-service/server"
)

// IsEnabled for gRPC.
func IsEnabled(cfg *Config) bool {
	return cfg != nil && server.IsEnabled(cfg.Config)
}

// UserAgent for gRPC.
func UserAgent(cfg *Config) string {
	if IsEnabled(cfg) {
		return server.UserAgent(cfg.Config)
	}

	return ""
}

// Config for gRPC.
type Config struct {
	*server.Config `yaml:",inline" json:",inline" toml:",inline"`
}
