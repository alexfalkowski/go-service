package client

import (
	"github.com/alexfalkowski/go-service/v2/config/options"
	"github.com/alexfalkowski/go-service/v2/crypto/tls"
	"github.com/alexfalkowski/go-service/v2/limiter"
	"github.com/alexfalkowski/go-service/v2/retry"
	"github.com/alexfalkowski/go-service/v2/token"
)

// Config configures client-side behavior shared across transports.
type Config struct {
	// Limiter configures client-side request limiting.
	Limiter *limiter.Config `yaml:"limiter,omitempty" json:"limiter,omitempty" toml:"limiter,omitempty"`

	// TLS configures client-side TLS settings (for example trusted roots, certs, and keys where applicable).
	TLS *tls.Config `yaml:"tls,omitempty" json:"tls,omitempty" toml:"tls,omitempty"`

	// Token configures client-side token handling (for example minting/attaching auth tokens).
	Token *token.Config `yaml:"token,omitempty" json:"token,omitempty" toml:"token,omitempty"`

	// Retry configures client-side retry behavior for outbound requests.
	Retry *retry.Config `yaml:"retry,omitempty" json:"retry,omitempty" toml:"retry,omitempty"`

	// Options provides client/transport-specific options as key-value pairs.
	Options options.Map `yaml:"options,omitempty" json:"options,omitempty" toml:"options,omitempty"`

	// Address is the remote address for the client to connect to (for example "host:port").
	Address string `yaml:"address,omitempty" json:"address,omitempty" toml:"address,omitempty"`

	// Timeout is the client request timeout duration, encoded as a Go duration string (for example "30s", "5m").
	Timeout string `yaml:"timeout,omitempty" json:"timeout,omitempty" toml:"timeout,omitempty"`
}

// IsEnabled reports whether client configuration is present.
func (c *Config) IsEnabled() bool {
	return c != nil
}
