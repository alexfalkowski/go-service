package hooks

import "github.com/alexfalkowski/go-service/v2/os"

// Config for hooks.
type Config struct {
	Secret string `yaml:"secret,omitempty" json:"secret,omitempty" toml:"secret,omitempty"`
}

// IsEnabled for hooks.
func (c *Config) IsEnabled() bool {
	return c != nil
}

// GetCert for hooks.
func (c *Config) GetSecret(fs *os.FS) ([]byte, error) {
	return fs.ReadSource(c.Secret)
}
