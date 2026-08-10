package http

import (
	"github.com/alexfalkowski/go-service/v2/config/server"
	"github.com/alexfalkowski/go-service/v2/time"
)

// Config configures the HTTP transport stack.
//
// It embeds the shared server-side transport configuration, which includes common fields such as:
//
//   - Address binding
//   - TLS configuration
//   - Low-level server options
//
// Timeout is a unary request execution budget. Socket and streaming deadlines
// are configured independently through the low-level options.
//
// This config is typically nested under the top-level `transport.Config` and is used by constructors
// such as [NewServer] to decide whether the HTTP transport should be wired.
//
// The struct tags are compatible with the repository's config decoder (YAML/JSON/TOML).
type Config struct {
	*server.Config `yaml:",inline" json:",inline" toml:",inline"`

	// Timeout bounds non-streaming request execution.
	//
	// In config files it is encoded as a Go duration string (for example "30s", "5m").
	// A zero value applies [time.DefaultTimeout]. Negative values are invalid.
	Timeout time.Duration `yaml:"timeout,omitempty" json:"timeout,omitempty" toml:"timeout,omitempty" validate:"gte=0"`
}

// IsEnabled reports whether the HTTP transport is enabled.
//
// It returns false when the receiver is nil, which allows config to be omitted entirely to disable the
// HTTP transport stack.
func (c *Config) IsEnabled() bool {
	return c != nil && c.Config.IsEnabled()
}

// GetTimeout returns the configured unary request execution timeout.
//
// A nil receiver or a non-positive value falls back to [time.DefaultTimeout].
func (c *Config) GetTimeout() time.Duration {
	if c == nil || c.Timeout <= 0 {
		return time.DefaultTimeout
	}

	return c.Timeout
}

// GetReadTimeout returns the configured streaming read timeout. The
// "read_timeout" option overrides the lower-level [time.DefaultTimeout].
func (c *Config) GetReadTimeout() time.Duration {
	if c == nil {
		return time.DefaultTimeout
	}

	return c.Options.NonNegativeDuration("read_timeout", time.DefaultTimeout)
}

// GetWriteTimeout returns the configured streaming write timeout. The
// "write_timeout" option overrides the lower-level [time.DefaultTimeout].
func (c *Config) GetWriteTimeout() time.Duration {
	if c == nil {
		return time.DefaultTimeout
	}

	return c.Options.NonNegativeDuration("write_timeout", time.DefaultTimeout)
}
