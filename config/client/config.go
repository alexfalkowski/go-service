package client

import (
	"github.com/alexfalkowski/go-service/v2/config/options"
	"github.com/alexfalkowski/go-service/v2/crypto/tls"
	tlsconfig "github.com/alexfalkowski/go-service/v2/crypto/tls/config"
	"github.com/alexfalkowski/go-service/v2/limiter"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/retry"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/token"
)

// Config configures client-side behavior shared across transports.
type Config struct {
	// Limiter configures client-side request limiting.
	Limiter *limiter.Config `yaml:"limiter,omitempty" json:"limiter,omitempty" toml:"limiter,omitempty"`

	// TLS configures client-side TLS settings (for example trusted roots, certs, and keys where applicable).
	TLS *tlsconfig.Config `yaml:"tls,omitempty" json:"tls,omitempty" toml:"tls,omitempty"`

	// Token configures client-side token handling (for example minting/attaching auth tokens).
	Token *token.Config `yaml:"token,omitempty" json:"token,omitempty" toml:"token,omitempty"`

	// Retry configures client-side retry behavior for outbound requests.
	Retry *retry.Config `yaml:"retry,omitempty" json:"retry,omitempty" toml:"retry,omitempty"`

	// Options provides client/transport-specific options as key-value pairs.
	Options options.Map `yaml:"options,omitempty" json:"options,omitempty" toml:"options,omitempty"`

	// Address is the remote address for the client to connect to (for example "host:port").
	Address string `yaml:"address,omitempty" json:"address,omitempty" toml:"address,omitempty"`

	// Timeout is the client request timeout duration.
	//
	// In config files it is encoded as a Go duration string (for example "30s", "5m").
	Timeout time.Duration `yaml:"timeout,omitempty" json:"timeout,omitempty" toml:"timeout,omitempty"`
}

// IsEnabled reports whether client configuration is present.
func (c *Config) IsEnabled() bool {
	return c != nil
}

// NewConfig constructs a client-side runtime TLS config from cfg.
//
// It uses CA as RootCAs and uses any configured certificate/key pair as the
// client certificate presented to servers that request one. If ServerName is
// configured, it is copied into the runtime TLS config for server certificate
// hostname verification.
func NewConfig(fs *os.FS, cfg *tlsconfig.Config) (*tls.Config, error) {
	config := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	if !cfg.IsEnabled() {
		return config, nil
	}

	if cfg.HasKeyPair() {
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

		config.RootCAs = pool
	}

	config.ServerName = cfg.ServerName

	return config, nil
}
