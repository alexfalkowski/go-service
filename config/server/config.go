package server

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/config/options"
	"github.com/alexfalkowski/go-service/v2/crypto/tls"
	tlsconfig "github.com/alexfalkowski/go-service/v2/crypto/tls/config"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/token"
	"github.com/alexfalkowski/go-service/v2/transport/limiter"
)

// ErrMissingKeyPair is returned when server TLS is configured without a complete certificate/key pair.
var ErrMissingKeyPair = errors.New("server: missing tls key pair")

// Config configures server-side behavior shared across transports.
type Config struct {
	// Limiter configures server-side request limiting.
	Limiter *limiter.Config `yaml:"limiter,omitempty" json:"limiter,omitempty" toml:"limiter,omitempty"`

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
	//
	// A zero value applies time.DefaultTimeout. Negative values are invalid.
	Timeout time.Duration `yaml:"timeout,omitempty" json:"timeout,omitempty" toml:"timeout,omitempty" validate:"gte=0"`

	// MaxReceiveSize limits inbound request payload size.
	//
	// In config files it is encoded as a human-readable SI size string (for example "64B", "2MB", "4GB").
	//
	// A zero value applies bytes.DefaultSize. Negative values are invalid.
	MaxReceiveSize bytes.Size `yaml:"max_receive_size,omitempty" json:"max_receive_size,omitempty" toml:"max_receive_size,omitempty" validate:"gte=0"`
}

// IsEnabled reports whether server configuration is present.
func (c *Config) IsEnabled() bool {
	return c != nil
}

// GetMaxReceiveSize returns the configured inbound payload limit.
//
// A nil receiver or a zero value falls back to bytes.DefaultSize.
func (c *Config) GetMaxReceiveSize() bytes.Size {
	if c == nil || c.MaxReceiveSize == 0 {
		return bytes.DefaultSize
	}

	return c.MaxReceiveSize
}

// GetTimeout returns the configured server timeout.
//
// A nil receiver or a non-positive value falls back to time.DefaultTimeout.
func (c *Config) GetTimeout() time.Duration {
	if c == nil || c.Timeout <= 0 {
		return time.DefaultTimeout
	}

	return c.Timeout
}

// NewConfig constructs a server-side runtime TLS config from cfg.
//
// If cfg has a CA configured, NewConfig uses it as ClientCAs and sets
// ClientAuth to tls.RequireAndVerifyClientCert. Without CA, the server config
// does not request client certificates.
//
// Any configured TLS material requires a complete server certificate/key pair.
// A CA-only configuration returns ErrMissingKeyPair.
func NewConfig(fs *os.FS, cfg *tlsconfig.Config) (*tls.Config, error) {
	config := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	if !cfg.HasKeyMaterial() && !cfg.HasCA() {
		return config, nil
	}

	if !cfg.HasKeyPair() {
		return config, ErrMissingKeyPair
	}

	pair, err := tlsconfig.NewKeyPair(fs, cfg)
	if err != nil {
		return config, err
	}

	config.Certificates = []tls.Certificate{pair}

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
