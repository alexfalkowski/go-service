package hmac

import "github.com/alexfalkowski/go-service/v2/os"

// Config configures HMAC key material loading for the HMAC-SHA-512 signer wired by this package.
type Config struct {
	// Key is a "source string" that resolves to the raw HMAC key bytes.
	//
	// It supports the go-service source string pattern implemented by `os.FS.ReadSource`:
	//   - "env:NAME" to read from an environment variable,
	//   - "file:/path/to/key" to read from a file, or
	//   - any other value treated as a literal.
	//
	// The resolved key bytes are used as the HMAC secret. Choose a high-entropy key and keep it private.
	Key string `yaml:"key,omitempty" json:"key,omitempty" toml:"key,omitempty"`
}

// IsEnabled reports whether HMAC configuration is enabled.
//
// By convention, a nil *Config is treated as "HMAC disabled" by wiring that depends on this configuration.
func (c *Config) IsEnabled() bool {
	return c != nil
}

// GetKey resolves and returns the HMAC key bytes using the configured Key source.
//
// It delegates to `fs.ReadSource(c.Key)` and returns any read/resolve error from that operation.
func (c *Config) GetKey(fs *os.FS) ([]byte, error) {
	return fs.ReadSource(c.Key)
}
