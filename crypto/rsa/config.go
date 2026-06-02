package rsa

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"

	"github.com/alexfalkowski/go-service/v2/crypto/errors"
	"github.com/alexfalkowski/go-service/v2/crypto/pem"
)

// Config configures RSA key loading for the RSA-OAEP encryptor/decryptor wired by this package.
//
// Public and Private are "source strings" resolved by [github.com/alexfalkowski/go-service/v2/crypto/pem.Decoder] (for example "env:NAME", "file:/path",
// or a literal PEM value).
//
// Expected key formats:
//   - Public: PEM block "RSA PUBLIC KEY" containing PKCS#1-encoded bytes (parsed via [crypto/x509.ParsePKCS1PublicKey]).
//   - Private: PEM block "RSA PRIVATE KEY" containing PKCS#1-encoded bytes (parsed via [crypto/x509.ParsePKCS1PrivateKey]).
type Config struct {
	// Public is a "source string" for the RSA public key PEM.
	//
	// It must decode to a PEM block of type "RSA PUBLIC KEY" with PKCS#1-encoded bytes.
	Public string `yaml:"public,omitempty" json:"public,omitempty" toml:"public,omitempty"`

	// Private is a "source string" for the RSA private key PEM.
	//
	// It must decode to a PEM block of type "RSA PRIVATE KEY" with PKCS#1-encoded bytes.
	Private string `yaml:"private,omitempty" json:"private,omitempty" toml:"private,omitempty"`
}

// IsEnabled reports whether RSA configuration is enabled.
//
// By convention, a nil *[Config] is treated as "RSA disabled" by wiring that depends on this configuration.
func (c *Config) IsEnabled() bool {
	return c != nil
}

// PublicKey loads and parses the configured RSA public key.
//
// It decodes a PEM "RSA PUBLIC KEY" block and parses it as a PKCS#1 public key ([crypto/x509.ParsePKCS1PublicKey]).
func (c *Config) PublicKey(decoder *pem.Decoder) (*rsa.PublicKey, error) {
	d, err := decoder.Decode(c.Public, "RSA PUBLIC KEY")
	if err != nil {
		return nil, err
	}

	key, err := x509.ParsePKCS1PublicKey(d)
	if err != nil {
		return nil, err
	}
	if err := validatePublicKey(key); err != nil {
		return nil, err
	}

	return key, nil
}

// PrivateKey loads and parses the configured RSA private key.
//
// It decodes a PEM "RSA PRIVATE KEY" block and parses it as a PKCS#1 private key ([crypto/x509.ParsePKCS1PrivateKey]).
func (c *Config) PrivateKey(decoder *pem.Decoder) (*rsa.PrivateKey, error) {
	d, err := decoder.Decode(c.Private, "RSA PRIVATE KEY")
	if err != nil {
		return nil, err
	}

	key, err := x509.ParsePKCS1PrivateKey(d)
	if err != nil {
		return nil, err
	}
	if err := validatePrivateKey(key); err != nil {
		return nil, err
	}

	return key, nil
}

func validatePublicKey(key *rsa.PublicKey) error {
	size := key.N.BitLen()
	if size < KeySize {
		return fmt.Errorf("rsa: invalid public key size %d: %w", size, errors.ErrInvalidKeySize)
	}

	return nil
}

func validatePrivateKey(key *rsa.PrivateKey) error {
	size := key.N.BitLen()
	if size < KeySize {
		return fmt.Errorf("rsa: invalid private key size %d: %w", size, errors.ErrInvalidKeySize)
	}

	return key.Validate()
}
