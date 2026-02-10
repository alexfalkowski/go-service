package grpc

import "github.com/alexfalkowski/go-service/v2/config/server"

// Config configures the gRPC transport.
//
// It embeds the shared server-side transport configuration.
type Config struct {
	*server.Config `yaml:",inline" json:",inline" toml:",inline"`
}

// IsEnabled reports whether the transport is enabled.
func (c *Config) IsEnabled() bool {
	return c != nil && c.Config.IsEnabled()
}
