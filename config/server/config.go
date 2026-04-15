package server

import (
	"github.com/alexfalkowski/go-service/v2/config/options"
	"github.com/alexfalkowski/go-service/v2/crypto/tls"
	"github.com/alexfalkowski/go-service/v2/limiter"
	"github.com/alexfalkowski/go-service/v2/retry"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/token"
)

// DefaultMaxReceiveBytes is the default inbound payload limit applied when MaxReceiveBytes is omitted or zero.
const DefaultMaxReceiveBytes int64 = 4 * 1024 * 1024

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

	// Address is the bind address for the server (for example ":8080").
	Address string `yaml:"address,omitempty" json:"address,omitempty" toml:"address,omitempty"`

	// Timeout is the server request timeout duration.
	//
	// In config files it is encoded as a Go duration string (for example "30s", "5m").
	Timeout time.Duration `yaml:"timeout,omitempty" json:"timeout,omitempty" toml:"timeout,omitempty"`

	// MaxReceiveBytes limits inbound request payload size in bytes.
	//
	// A zero value applies DefaultMaxReceiveBytes. Negative values are invalid.
	MaxReceiveBytes int64 `yaml:"max_receive_bytes,omitempty" json:"max_receive_bytes,omitempty" toml:"max_receive_bytes,omitempty" validate:"gte=0"`
}

// IsEnabled reports whether server configuration is present.
func (c *Config) IsEnabled() bool {
	return c != nil
}

// GetMaxReceiveBytes returns the configured inbound payload limit in bytes.
//
// A nil receiver or a zero value falls back to DefaultMaxReceiveBytes.
func (c *Config) GetMaxReceiveBytes() int64 {
	if c == nil || c.MaxReceiveBytes == 0 {
		return DefaultMaxReceiveBytes
	}

	return c.MaxReceiveBytes
}
