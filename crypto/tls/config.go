package tls

import (
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// Config configures TLS key material loading.
//
// Cert and Key are "source strings" resolved by `os.FS.ReadSource`.
// They may be:
//   - "env:NAME" to read PEM bytes from environment variable NAME,
//   - "file:/path/to/pem" to read PEM bytes from a file, or
//   - any other value treated as the literal PEM content.
//
// This config is intentionally minimal: it only models a leaf certificate/key pair and does not
// represent trust roots, client CA pools, cipher suites, or other `crypto/tls.Config` knobs.
type Config struct {
	// Cert is a "source string" for the TLS certificate (PEM-encoded).
	//
	// The resolved value must contain a PEM-encoded certificate suitable for tls.X509KeyPair.
	Cert string `yaml:"cert,omitempty" json:"cert,omitempty" toml:"cert,omitempty"`

	// Key is a "source string" for the TLS private key (PEM-encoded).
	//
	// The resolved value must contain a PEM-encoded private key suitable for tls.X509KeyPair.
	Key string `yaml:"key,omitempty" json:"key,omitempty" toml:"key,omitempty"`
}

// IsEnabled reports whether TLS configuration is present.
//
// By convention across go-service config types, a nil *Config is treated as "disabled".
func (c *Config) IsEnabled() bool {
	return c != nil
}

// HasKeyPair reports whether both certificate and key sources are configured.
//
// This only checks that both source strings are non-empty; it does not validate that the resolved
// contents are well-formed PEM or that they form a valid X.509 key pair.
func (c *Config) HasKeyPair() bool {
	return !strings.IsEmpty(c.Cert) && !strings.IsEmpty(c.Key)
}

// GetCert resolves and returns the certificate bytes from the configured source string.
//
// It delegates to `fs.ReadSource(c.Cert)` and returns any read/resolve error from that operation.
func (c *Config) GetCert(fs *os.FS) ([]byte, error) {
	return fs.ReadSource(c.Cert)
}

// GetKey resolves and returns the private key bytes from the configured source string.
//
// It delegates to `fs.ReadSource(c.Key)` and returns any read/resolve error from that operation.
func (c *Config) GetKey(fs *os.FS) ([]byte, error) {
	return fs.ReadSource(c.Key)
}
