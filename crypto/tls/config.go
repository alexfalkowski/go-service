package tls

import (
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// Config configures TLS key material loading.
//
// Cert and Key are "source strings" that may be literal values, `file:` paths, or `env:` references.
type Config struct {
	Cert string `yaml:"cert,omitempty" json:"cert,omitempty" toml:"cert,omitempty"`
	Key  string `yaml:"key,omitempty" json:"key,omitempty" toml:"key,omitempty"`
}

// IsEnabled reports whether TLS is enabled (i.e., the config is non-nil).
func (c *Config) IsEnabled() bool {
	return c != nil
}

// HasKeyPair reports whether both certificate and key sources are configured.
func (c *Config) HasKeyPair() bool {
	return !strings.IsEmpty(c.Cert) && !strings.IsEmpty(c.Key)
}

// GetCert reads the certificate bytes from the configured source.
func (c *Config) GetCert(fs *os.FS) ([]byte, error) {
	return fs.ReadSource(c.Cert)
}

// GetKey reads the private key bytes from the configured source.
func (c *Config) GetKey(fs *os.FS) ([]byte, error) {
	return fs.ReadSource(c.Key)
}
