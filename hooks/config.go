package hooks

import "github.com/alexfalkowski/go-service/v2/os"

// Config configures Standard Webhooks secret loading.
type Config struct {
	// Secret is a source string for the webhook secret.
	//
	// It supports the go-service "source string" pattern (for example "env:NAME", "file:/path", or a literal value),
	// as implemented by os.FS.ReadSource.
	Secret string `yaml:"secret,omitempty" json:"secret,omitempty" toml:"secret,omitempty"`
}

// IsEnabled reports whether hooks are enabled.
func (c *Config) IsEnabled() bool {
	return c != nil
}

// GetSecret loads and returns the webhook secret bytes using the configured Secret source.
func (c *Config) GetSecret(fs *os.FS) ([]byte, error) {
	return fs.ReadSource(c.Secret)
}
