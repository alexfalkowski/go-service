package aes

import "github.com/alexfalkowski/go-service/v2/os"

// Config configures AES-GCM key loading.
type Config struct {
	// Key is a source string for the AES key material.
	//
	// It supports the go-service "source string" pattern (for example "env:NAME", "file:/path", or a literal value),
	// as implemented by os.FS.ReadSource.
	Key string `yaml:"key,omitempty" json:"key,omitempty" toml:"key,omitempty"`
}

// IsEnabled reports whether AES configuration is enabled.
//
// A nil config is considered disabled.
func (c *Config) IsEnabled() bool {
	return c != nil
}

// GetKey loads and returns the AES key bytes using the configured Key source.
func (c *Config) GetKey(fs *os.FS) ([]byte, error) {
	return fs.ReadSource(c.Key)
}
