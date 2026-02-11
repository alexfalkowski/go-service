package server

import (
	"github.com/alexfalkowski/go-service/v2/config/options"
	"github.com/alexfalkowski/go-service/v2/crypto/tls"
	"github.com/alexfalkowski/go-service/v2/limiter"
	"github.com/alexfalkowski/go-service/v2/retry"
	"github.com/alexfalkowski/go-service/v2/token"
)

// Config configures server-side behavior shared across transports.
type Config struct {
	// Limiter configures server-side request limiting.
	Limiter *limiter.Config `yaml:"limiter,omitempty" json:"limiter,omitempty" toml:"limiter,omitempty"`

	// Retry configures server-side retry behavior where applicable.
	Retry *retry.Config `yaml:"retry,omitempty" json:"retry,omitempty" toml:"retry,omitempty"`

	// TLS configures server-side TLS (certificate and key material).
	TLS *tls.Config `yaml:"tls,omitempty" json:"tls,omitempty" toml:"tls,omitempty"`

	// Token configures server-side token validation/handling.
	Token *token.Config `yaml:"token,omitempty" json:"token,omitempty" toml:"token,omitempty"`

	// Options provides transport/server-specific options as key-value pairs.
	Options options.Map `yaml:"options,omitempty" json:"options,omitempty" toml:"options,omitempty"`

	// Timeout is the server request timeout duration, encoded as a Go duration string (for example "30s", "5m").
	Timeout string `yaml:"timeout,omitempty" json:"timeout,omitempty" toml:"timeout,omitempty"`

	// Address is the bind address for the server (for example ":8080").
	Address string `yaml:"address,omitempty" json:"address,omitempty" toml:"address,omitempty"`
}

// IsEnabled for server.
func (c *Config) IsEnabled() bool {
	return c != nil
}
