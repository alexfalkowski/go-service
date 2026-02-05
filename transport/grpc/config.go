package grpc

import "github.com/alexfalkowski/go-service/v2/config/server"

// Config for gRPC.
type Config struct {
	*server.Config `yaml:",inline" json:",inline" toml:",inline"`
}

// IsEnabled for transport.
func (c *Config) IsEnabled() bool {
	return c != nil && c.Config.IsEnabled()
}
