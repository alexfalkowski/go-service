package aes

import "github.com/alexfalkowski/go-service/v2/os"

// Config for aes.
type Config struct {
	Key string `yaml:"key,omitempty" json:"key,omitempty" toml:"key,omitempty"`
}

// IsEnabled for aes.
func (c *Config) IsEnabled() bool {
	return c != nil
}

// GetKey for aes.
func (c *Config) GetKey(fs *os.FS) ([]byte, error) {
	return fs.ReadSource(c.Key)
}
