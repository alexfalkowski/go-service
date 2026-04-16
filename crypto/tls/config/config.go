package config

import (
	"github.com/alexfalkowski/go-service/v2/crypto/tls"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// Config configures TLS key material loading from go-service source strings.
//
// Cert and Key are "source strings" resolved by `os.FS.ReadSource`.
// They may be:
//   - "env:NAME" to read PEM bytes from environment variable NAME,
//   - "file:/path/to/pem" to read PEM bytes from a file, or
//   - any other value treated as the literal PEM content.
//
// This config is intentionally minimal: it only models the leaf
// certificate/private-key pair that `NewConfig` loads into a runtime TLS
// config. It does not model trust roots, client CA pools, cipher suites,
// ALPN, session tickets, or the many other knobs on `crypto/tls.Config`.
type Config struct {
	// Cert is a "source string" for the TLS certificate (PEM-encoded).
	//
	// The resolved value must contain a PEM-encoded certificate suitable for
	// tls.X509KeyPair. Its contents are not parsed or validated until
	// `NewConfig` is called.
	Cert string `yaml:"cert,omitempty" json:"cert,omitempty" toml:"cert,omitempty"`

	// Key is a "source string" for the TLS private key (PEM-encoded).
	//
	// The resolved value must contain a PEM-encoded private key suitable for
	// tls.X509KeyPair. Its contents are not parsed or validated until
	// `NewConfig` is called.
	Key string `yaml:"key,omitempty" json:"key,omitempty" toml:"key,omitempty"`
}

// IsEnabled reports whether TLS configuration is present.
//
// By convention across go-service config types, a nil *Config is treated as
// "disabled", so callers can omit TLS config entirely to leave transports in
// plain-text mode.
func (c *Config) IsEnabled() bool {
	return c != nil
}

// HasKeyPair reports whether both certificate and key sources are configured.
//
// This only checks that both source strings are non-empty. It does not validate
// that the resolved contents are readable, well-formed PEM, or that they form a
// valid X.509 key pair.
func (c *Config) HasKeyPair() bool {
	return !strings.IsEmpty(c.Cert) && !strings.IsEmpty(c.Key)
}

// GetCert resolves and returns the certificate bytes from the configured source
// string.
//
// It delegates to `fs.ReadSource(c.Cert)` and returns any read/resolve error
// from that operation.
func (c *Config) GetCert(fs *os.FS) ([]byte, error) {
	return fs.ReadSource(c.Cert)
}

// GetKey resolves and returns the private key bytes from the configured source
// string.
//
// It delegates to `fs.ReadSource(c.Key)` and returns any read/resolve error
// from that operation.
func (c *Config) GetKey(fs *os.FS) ([]byte, error) {
	return fs.ReadSource(c.Key)
}

// NewConfig constructs a runtime `*crypto/tls.Config` from cfg.
//
// # Defaults
//
// NewConfig always applies the following defaults:
//   - MinVersion: TLS 1.2 (tls.VersionTLS12)
//   - ClientAuth: require and verify client certificates (mTLS) via tls.RequireAndVerifyClientCert
//
// These defaults are conservative service-to-service defaults. Callers that
// need to adjust additional runtime TLS settings can modify the returned config
// after construction.
//
// # Key material loading
//
// If cfg is nil (disabled) or cfg does not have both a certificate and key configured (see cfg.HasKeyPair),
// NewConfig returns a config with the defaults above and does not attempt to load any key material.
//
// When cfg is enabled and has a key pair, certificate and key bytes are resolved via the provided filesystem
// using go-service "source strings" (literal value, `file:` path, or `env:` reference). The resolved PEM
// is then parsed using tls.X509KeyPair and attached to the returned config.
//
// # Errors
//
// NewConfig returns the partially constructed config along with any error encountered while resolving the
// certificate/key sources or parsing them as an X.509 key pair.
//
// The returned config contains no leaf certificates when TLS is disabled or no
// key pair is configured.
func NewConfig(fs *os.FS, cfg *Config) (*tls.Config, error) {
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}

	if !cfg.IsEnabled() || !cfg.HasKeyPair() {
		return tlsConfig, nil
	}

	cert, err := cfg.GetCert(fs)
	if err != nil {
		return tlsConfig, err
	}

	key, err := cfg.GetKey(fs)
	if err != nil {
		return tlsConfig, err
	}

	pair, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return tlsConfig, err
	}

	tlsConfig.Certificates = []tls.Certificate{pair}
	return tlsConfig, nil
}
