package ed25519

import (
	"crypto/ed25519"
	"crypto/x509"

	"github.com/alexfalkowski/go-service/v2/crypto/pem"
)

// Config configures Ed25519 key loading.
//
// Public and Private are "source strings" that are read by crypto/pem.Decoder (for example "env:NAME", "file:/path",
// or a literal PEM value).
//
// Expected key formats:
//   - Public: PEM block "PUBLIC KEY" containing PKIX-encoded bytes (x509.ParsePKIXPublicKey).
//   - Private: PEM block "PRIVATE KEY" containing PKCS#8-encoded bytes (x509.ParsePKCS8PrivateKey).
//
// Note: the key parsing helpers on this type use type assertions and will panic if the decoded key is not the
// expected Ed25519 key type. See PublicKey and PrivateKey for details.
type Config struct {
	// Public is a "source string" for the Ed25519 public key PEM.
	//
	// It is decoded by crypto/pem.Decoder and must contain a PEM block of type "PUBLIC KEY"
	// with PKIX-encoded bytes (parsed via x509.ParsePKIXPublicKey).
	Public string `yaml:"public,omitempty" json:"public,omitempty" toml:"public,omitempty"`

	// Private is a "source string" for the Ed25519 private key PEM.
	//
	// It is decoded by crypto/pem.Decoder and must contain a PEM block of type "PRIVATE KEY"
	// with PKCS#8-encoded bytes (parsed via x509.ParsePKCS8PrivateKey).
	Private string `yaml:"private,omitempty" json:"private,omitempty" toml:"private,omitempty"`
}

// IsEnabled reports whether Ed25519 configuration is enabled.
//
// A nil config is considered disabled.
func (c *Config) IsEnabled() bool {
	return c != nil
}

// PublicKey loads and parses the configured Ed25519 public key.
//
// It decodes a PEM "PUBLIC KEY" block from c.Public and parses it as a PKIX public key
// (x509.ParsePKIXPublicKey).
//
// Panics: if the decoded key is not an ed25519.PublicKey, this function will panic due to a type assertion.
// This can happen if the PEM data is a valid "PUBLIC KEY" block but contains a different key type
// (for example RSA or ECDSA).
func (c *Config) PublicKey(decoder *pem.Decoder) (ed25519.PublicKey, error) {
	d, err := decoder.Decode(c.Public, "PUBLIC KEY")
	if err != nil {
		return nil, err
	}

	k, err := x509.ParsePKIXPublicKey(d)
	if err != nil {
		return nil, err
	}

	return k.(ed25519.PublicKey), nil
}

// PrivateKey loads and parses the configured Ed25519 private key.
//
// It decodes a PEM "PRIVATE KEY" block from c.Private and parses it as a PKCS#8 private key
// (x509.ParsePKCS8PrivateKey).
//
// Panics: if the decoded key is not an ed25519.PrivateKey, this function will panic due to a type assertion.
// This can happen if the PEM data is a valid "PRIVATE KEY" block but contains a different key type.
func (c *Config) PrivateKey(decoder *pem.Decoder) (ed25519.PrivateKey, error) {
	d, err := decoder.Decode(c.Private, "PRIVATE KEY")
	if err != nil {
		return nil, err
	}

	k, err := x509.ParsePKCS8PrivateKey(d)
	if err != nil {
		return nil, err
	}

	return k.(ed25519.PrivateKey), nil
}
