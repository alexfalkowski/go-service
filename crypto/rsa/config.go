package rsa

import (
	"crypto/rsa"
	"crypto/x509"

	"github.com/alexfalkowski/go-service/crypto/pem"
)

// IsEnabled for rsa.
func IsEnabled(cfg *Config) bool {
	return cfg != nil
}

type (
	// PublicKey for rsa.
	PublicKey string

	// PrivateKey for rsa.
	PrivateKey string

	// Config for rsa.
	Config struct {
		Public  PublicKey  `yaml:"public,omitempty" json:"public,omitempty" toml:"public,omitempty"`
		Private PrivateKey `yaml:"private,omitempty" json:"private,omitempty" toml:"private,omitempty"`
	}
)

// PublicKey rsa.
func (c *Config) PublicKey() (*rsa.PublicKey, error) {
	d, err := pem.Decode(string(c.Public), "RSA PUBLIC KEY")
	if err != nil {
		return nil, err
	}

	return x509.ParsePKCS1PublicKey(d)
}

// PrivateKey rsa.
func (c *Config) PrivateKey() (*rsa.PrivateKey, error) {
	d, err := pem.Decode(string(c.Private), "RSA PRIVATE KEY")
	if err != nil {
		return nil, err
	}

	return x509.ParsePKCS1PrivateKey(d)
}
