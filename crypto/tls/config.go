package tls

import (
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// IsEnabled for security.
func IsEnabled(cfg *Config) bool {
	return cfg != nil
}

// Config for tls.
type Config struct {
	Cert string `yaml:"cert,omitempty" json:"cert,omitempty" toml:"cert,omitempty"`
	Key  string `yaml:"key,omitempty" json:"key,omitempty" toml:"key,omitempty"`
}

// HasKeyPair for security.
func (c *Config) HasKeyPair() bool {
	return !strings.IsEmpty(c.Cert) && !strings.IsEmpty(c.Key)
}

// GetCert for tls.
func (c *Config) GetCert(fs *os.FS) ([]byte, error) {
	return fs.ReadSource(c.Cert)
}

// GetKey for tls.
func (c *Config) GetKey(fs *os.FS) ([]byte, error) {
	return fs.ReadSource(c.Key)
}
