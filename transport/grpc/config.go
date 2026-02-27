package grpc

import "github.com/alexfalkowski/go-service/v2/config/server"

// Config configures the gRPC transport stack.
//
// It embeds the shared server-side transport configuration, which includes common fields such as:
//
//   - Address binding
//   - Server timeouts
//   - TLS configuration
//   - Low-level server options
//
// This config is typically nested under the top-level `transport.Config` and is used by constructors
// such as `NewServer` to decide whether the gRPC transport should be wired.
//
// The struct tags are compatible with the repository's config decoder (YAML/JSON/TOML).
type Config struct {
	*server.Config `yaml:",inline" json:",inline" toml:",inline"`
}

// IsEnabled reports whether the gRPC transport is enabled.
//
// It returns false when the receiver is nil, which allows config to be omitted entirely to disable the
// gRPC transport stack.
func (c *Config) IsEnabled() bool {
	return c != nil && c.Config.IsEnabled()
}
