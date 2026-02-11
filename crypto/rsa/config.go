package rsa

import (
	"crypto/rsa"
	"crypto/x509"

	"github.com/alexfalkowski/go-service/v2/crypto/pem"
)

// Config configures RSA key loading.
//
// Public and Private are "source strings" that are read by crypto/pem.Decoder (for example "env:NAME", "file:/path",
// or a literal PEM value).
//
// Expected key formats:
//   - Public: PEM block "RSA PUBLIC KEY" containing PKCS#1-encoded bytes (x509.ParsePKCS1PublicKey).
//   - Private: PEM block "RSA PRIVATE KEY" containing PKCS#1-encoded bytes (x509.ParsePKCS1PrivateKey).
type Config struct {
	// Public is a "source string" for the RSA public key PEM.
	//
	// It is decoded by crypto/pem.Decoder and must contain a PEM block of type "RSA PUBLIC KEY"
	// with PKCS#1-encoded bytes (parsed via x509.ParsePKCS1PublicKey).
	Public string `yaml:"public,omitempty" json:"public,omitempty" toml:"public,omitempty"`

	// Private is a "source string" for the RSA private key PEM.
	//
	// It is decoded by crypto/pem.Decoder and must contain a PEM block of type "RSA PRIVATE KEY"
	// with PKCS#1-encoded bytes (parsed via x509.ParsePKCS1PrivateKey).
	Private string `yaml:"private,omitempty" json:"private,omitempty" toml:"private,omitempty"`
}

// IsEnabled reports whether RSA configuration is enabled.
//
// A nil config is considered disabled.
func (c *Config) IsEnabled() bool {
	return c != nil
}

// PublicKey loads and parses the configured RSA public key.
//
// It decodes a PEM "RSA PUBLIC KEY" block and parses it as a PKCS#1 public key.
func (c *Config) PublicKey(decoder *pem.Decoder) (*rsa.PublicKey, error) {
	d, err := decoder.Decode(c.Public, "RSA PUBLIC KEY")
	if err != nil {
		return nil, err
	}

	return x509.ParsePKCS1PublicKey(d)
}

// PrivateKey loads and parses the configured RSA private key.
//
// It decodes a PEM "RSA PRIVATE KEY" block and parses it as a PKCS#1 private key.
func (c *Config) PrivateKey(decoder *pem.Decoder) (*rsa.PrivateKey, error) {
	d, err := decoder.Decode(c.Private, "RSA PRIVATE KEY")
	if err != nil {
		return nil, err
	}

	return x509.ParsePKCS1PrivateKey(d)
}
