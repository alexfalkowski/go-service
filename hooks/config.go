package hooks

import "github.com/alexfalkowski/go-service/v2/os"

// IsEnabled for hooks.
func IsEnabled(cfg *Config) bool {
	return cfg != nil
}

// Config for hooks.
type Config struct {
	Secret string `yaml:"secret,omitempty" json:"secret,omitempty" toml:"secret,omitempty"`
}

// GetCert for hooks.
func (c *Config) GetSecret(fs *os.FS) ([]byte, error) {
	return fs.ReadSource(c.Secret)
}
