package debug

import "github.com/alexfalkowski/go-service/v2/config/server"

// Config configures the debug server.
//
// It embeds [github.com/alexfalkowski/go-service/v2/config/server.Config] to reuse common server settings such as address,
// TLS configuration, maximum inbound request-body size, and server-specific options.
//
// # Optional pointers and "enabled" semantics
//
// This config is intentionally optional. By convention across go-service configuration types, a nil
// *[Config] is treated as "debug disabled". A present config is enabled even when none of the embedded
// [server.Config] fields are set, in which case the shared defaults apply.
//
// This allows services to omit the debug config entirely to disable the debug server.
type Config struct {
	server.Config `yaml:",inline" json:",inline" toml:",inline"`
}

// IsEnabled reports whether debug configuration is present.
//
// It returns false when the receiver is nil, which allows config to be omitted entirely to disable the
// debug server.
func (c *Config) IsEnabled() bool {
	return c != nil
}
