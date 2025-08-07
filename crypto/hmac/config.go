package hmac

import "github.com/alexfalkowski/go-service/v2/os"

// Config for hmac.
type Config struct {
	Key string `yaml:"key,omitempty" json:"key,omitempty" toml:"key,omitempty"`
}

// IsEnabled for hmac.
func (c *Config) IsEnabled() bool {
	return c != nil
}

// GetKey for hmac.
func (c *Config) GetKey(fs *os.FS) ([]byte, error) {
	return fs.ReadSource(c.Key)
}
