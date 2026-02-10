package hmac

import "github.com/alexfalkowski/go-service/v2/os"

// Config configures HMAC key loading.
type Config struct {
	// Key is a source string for the HMAC key material.
	//
	// It supports the go-service "source string" pattern (for example "env:NAME", "file:/path", or a literal value),
	// as implemented by os.FS.ReadSource.
	Key string `yaml:"key,omitempty" json:"key,omitempty" toml:"key,omitempty"`
}

// IsEnabled reports whether HMAC configuration is enabled.
//
// A nil config is considered disabled.
func (c *Config) IsEnabled() bool {
	return c != nil
}

// GetKey loads and returns the HMAC key bytes using the configured Key source.
func (c *Config) GetKey(fs *os.FS) ([]byte, error) {
	return fs.ReadSource(c.Key)
}
