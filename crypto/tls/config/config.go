package config

import (
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// ErrInvalidCA is returned when configured CA PEM does not contain any certificates.
var ErrInvalidCA = errors.New("tls: invalid ca")

// Config configures TLS key material loading from go-service source strings.
//
// Cert, Key, and CA are "source strings" resolved by `os.FS.ReadSource`.
// They may be:
//   - "env:NAME" to read PEM bytes from environment variable NAME,
//   - "file:/path/to/pem" to read PEM bytes from a file, or
//   - any other value treated as the literal PEM content.
//
// This config is intentionally minimal: it models leaf certificate/private-key
// material, an optional peer CA bundle, and an optional client-side server name.
// It does not model cipher suites, ALPN, session tickets, or the many other
// knobs on `crypto/tls.Config`.
type Config struct {
	// Cert is a "source string" for the TLS certificate (PEM-encoded).
	//
	// The resolved value must contain a PEM-encoded certificate suitable for
	// tls.X509KeyPair. Its contents are not parsed or validated until
	// a runtime TLS config is constructed.
	Cert string `yaml:"cert,omitempty" json:"cert,omitempty" toml:"cert,omitempty"`

	// Key is a "source string" for the TLS private key (PEM-encoded).
	//
	// The resolved value must contain a PEM-encoded private key suitable for
	// tls.X509KeyPair. Its contents are not parsed or validated until
	// a runtime TLS config is constructed.
	Key string `yaml:"key,omitempty" json:"key,omitempty" toml:"key,omitempty"`

	// CA is a "source string" for the peer CA bundle (PEM-encoded).
	//
	// Server TLS config uses CA as the client CA pool. Client TLS config uses CA
	// as the root CA pool for verifying the server certificate.
	CA string `yaml:"ca,omitempty" json:"ca,omitempty" toml:"ca,omitempty"`

	// ServerName is the optional server name clients use for certificate verification.
	//
	// Leave this empty when the transport can infer the server name from the
	// dial target or request URL. Set it when the dial address differs from the
	// certificate's DNS name, such as dialing 127.0.0.1 with a localhost cert.
	ServerName string `yaml:"server_name,omitempty" json:"server_name,omitempty" toml:"server_name,omitempty"`
}

// IsEnabled reports whether TLS configuration has any configured TLS field.
//
// By convention, nil and empty configs are disabled.
func (c *Config) IsEnabled() bool {
	return c.HasKeyMaterial() || c.HasCA() || c.HasServerName()
}

// HasKeyPair reports whether both certificate and key sources are configured.
//
// This only checks that both source strings are non-empty. It does not validate
// that the resolved contents are readable, well-formed PEM, or that they form a
// valid X.509 key pair.
func (c *Config) HasKeyPair() bool {
	return c != nil && !strings.IsEmpty(c.Cert) && !strings.IsEmpty(c.Key)
}

// HasKeyMaterial reports whether certificate or private key source is configured.
func (c *Config) HasKeyMaterial() bool {
	return c != nil && (!strings.IsEmpty(c.Cert) || !strings.IsEmpty(c.Key))
}

// HasCA reports whether a peer CA source is configured.
//
// This only checks that the source string is non-empty. It does not validate
// that the resolved contents are readable or contain PEM-encoded certificates.
func (c *Config) HasCA() bool {
	return c != nil && !strings.IsEmpty(c.CA)
}

// HasServerName reports whether a client-side server name override is configured.
func (c *Config) HasServerName() bool {
	return c != nil && !strings.IsEmpty(c.ServerName)
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

// GetCA resolves and returns the peer CA bytes from the configured source
// string.
//
// It delegates to `fs.ReadSource(c.CA)` and returns any read/resolve error
// from that operation.
func (c *Config) GetCA(fs *os.FS) ([]byte, error) {
	return fs.ReadSource(c.CA)
}
