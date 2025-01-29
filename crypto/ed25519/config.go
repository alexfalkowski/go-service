package ed25519

import (
	"crypto/ed25519"
	"crypto/x509"

	"github.com/alexfalkowski/go-service/crypto/pem"
	"github.com/alexfalkowski/go-service/structs"
)

// IsEnabled for ed25519.
func IsEnabled(cfg *Config) bool {
	return !structs.IsZero(cfg)
}

// Config for ed25519.
type Config struct {
	Public  string `yaml:"public,omitempty" json:"public,omitempty" toml:"public,omitempty"`
	Private string `yaml:"private,omitempty" json:"private,omitempty" toml:"private,omitempty"`
}

// PublicKey ed25519.
func (c *Config) PublicKey() (ed25519.PublicKey, error) {
	d, err := pem.Decode(c.Public, "PUBLIC KEY")
	if err != nil {
		return nil, err
	}

	k, err := x509.ParsePKIXPublicKey(d)
	if err != nil {
		return nil, err
	}

	return k.(ed25519.PublicKey), nil
}

// PrivateKey ed25519.
func (c *Config) PrivateKey() (ed25519.PrivateKey, error) {
	d, err := pem.Decode(c.Private, "PRIVATE KEY")
	if err != nil {
		return nil, err
	}

	k, err := x509.ParsePKCS8PrivateKey(d)
	if err != nil {
		return nil, err
	}

	return k.(ed25519.PrivateKey), nil
}
