package grpc

import (
	"github.com/alexfalkowski/go-service/server"
)

// IsEnabled for gRPC.
func IsEnabled(c *Config) bool {
	return c != nil && server.IsEnabled(c.Config)
}

// UserAgent for gRPC.
func UserAgent(c *Config) string {
	if IsEnabled(c) {
		return server.UserAgent(c.Config)
	}

	return ""
}

// Config for gRPC.
type Config struct {
	*server.Config `yaml:",inline" json:",inline" toml:",inline"`
}
