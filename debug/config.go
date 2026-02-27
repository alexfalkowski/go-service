package debug

import "github.com/alexfalkowski/go-service/v2/config/server"

// Config configures the debug server.
//
// It embeds `config/server.Config` to reuse common server settings such as address, timeout,
// TLS configuration, and server-specific options.
//
// # Optional pointers and "enabled" semantics
//
// This config is intentionally optional. By convention across go-service configuration types, a nil
// *Config is treated as "debug disabled". The embedded `*server.Config` is also optional; IsEnabled
// returns true only when both the outer *Config and the embedded *server.Config are non-nil/enabled.
//
// This allows services to omit the debug config entirely to disable the debug server.
type Config struct {
	*server.Config `yaml:",inline" json:",inline" toml:",inline"`
}

// IsEnabled reports whether debug configuration is present and enabled.
//
// It returns true only when both the debug wrapper config and the embedded server config are non-nil
// and enabled.
func (c *Config) IsEnabled() bool {
	return c != nil && c.Config.IsEnabled()
}
