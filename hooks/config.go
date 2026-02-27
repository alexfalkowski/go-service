package hooks

import "github.com/alexfalkowski/go-service/v2/os"

// Config configures Standard Webhooks secret loading.
type Config struct {
	// Secret is a "source string" that resolves to the webhook signing/verification secret bytes.
	//
	// It supports the go-service source string pattern implemented by `os.FS.ReadSource`:
	//   - "env:NAME" to read the secret from an environment variable,
	//   - "file:/path/to/secret" to read the secret from a file, or
	//   - any other value treated as the literal secret.
	//
	// Security note: keep this secret private and avoid logging it.
	Secret string `yaml:"secret,omitempty" json:"secret,omitempty" toml:"secret,omitempty"`
}

// IsEnabled reports whether hooks configuration is present.
//
// By convention across go-service config types, a nil *Config is treated as "disabled".
func (c *Config) IsEnabled() bool {
	return c != nil
}

// GetSecret resolves and returns the webhook secret bytes using the configured Secret source.
//
// It delegates to `fs.ReadSource(c.Secret)` and returns any read/resolve error from that operation.
func (c *Config) GetSecret(fs *os.FS) ([]byte, error) {
	return fs.ReadSource(c.Secret)
}
