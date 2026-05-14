package server

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/config/options"
	"github.com/alexfalkowski/go-service/v2/crypto/tls"
	tlsconfig "github.com/alexfalkowski/go-service/v2/crypto/tls/config"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/token"
	"github.com/alexfalkowski/go-service/v2/transport/limiter"
	"github.com/alexfalkowski/go-service/v2/transport/retry"
)

// DefaultMaxReceiveSize is the default inbound payload limit applied when MaxReceiveSize is omitted or zero.
const DefaultMaxReceiveSize bytes.Size = 4 * bytes.MB

// Config configures server-side behavior shared across transports.
type Config struct {
	// Limiter configures server-side request limiting.
	Limiter *limiter.Config `yaml:"limiter,omitempty" json:"limiter,omitempty" toml:"limiter,omitempty"`

	// Retry configures server-side retry behavior where applicable.
	Retry *retry.Config `yaml:"retry,omitempty" json:"retry,omitempty" toml:"retry,omitempty"`

	// TLS configures server-side TLS (certificate and key material).
	TLS *tlsconfig.Config `yaml:"tls,omitempty" json:"tls,omitempty" toml:"tls,omitempty"`

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

	// MaxReceiveSize limits inbound request payload size.
	//
	// In config files it is encoded as a human-readable SI size string (for example "64B", "2MB", "4GB").
	//
	// A zero value applies DefaultMaxReceiveSize. Negative values are invalid.
	MaxReceiveSize bytes.Size `yaml:"max_receive_size,omitempty" json:"max_receive_size,omitempty" toml:"max_receive_size,omitempty" validate:"gte=0"`
}

// IsEnabled reports whether server configuration is present.
func (c *Config) IsEnabled() bool {
	return c != nil
}

// GetMaxReceiveSize returns the configured inbound payload limit.
//
// A nil receiver or a zero value falls back to DefaultMaxReceiveSize.
func (c *Config) GetMaxReceiveSize() bytes.Size {
	if c == nil || c.MaxReceiveSize == 0 {
		return DefaultMaxReceiveSize
	}

	return c.MaxReceiveSize
}

// NewConfig constructs a server-side runtime TLS config from cfg.
//
// If cfg has a CA configured, NewConfig uses it as ClientCAs and sets
// ClientAuth to tls.RequireAndVerifyClientCert. Without CA, the server config
// does not request client certificates.
func NewConfig(fs *os.FS, cfg *tlsconfig.Config) (*tls.Config, error) {
	config := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	if !cfg.HasKeyMaterial() && !cfg.HasCA() {
		return config, nil
	}

	if cfg.HasKeyMaterial() {
		pair, err := tlsconfig.NewKeyPair(fs, cfg)
		if err != nil {
			return config, err
		}

		config.Certificates = []tls.Certificate{pair}
	}

	if cfg.HasCA() {
		pool, err := tlsconfig.NewCertPool(fs, cfg)
		if err != nil {
			return config, err
		}

		config.ClientCAs = pool
		config.ClientAuth = tls.RequireAndVerifyClientCert
	}

	return config, nil
}
