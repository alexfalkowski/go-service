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

// Config for rsa.
type Config struct {
	Public  string `yaml:"public,omitempty" json:"public,omitempty" toml:"public,omitempty"`
	Private string `yaml:"private,omitempty" json:"private,omitempty" toml:"private,omitempty"`
}

// PublicKey rsa.
func (c *Config) PublicKey() (*rsa.PublicKey, error) {
	d, err := pem.Decode(c.Public, "RSA PUBLIC KEY")
	if err != nil {
		return nil, err
	}

	return x509.ParsePKCS1PublicKey(d)
}

// PrivateKey rsa.
func (c *Config) PrivateKey() (*rsa.PrivateKey, error) {
	d, err := pem.Decode(c.Private, "RSA PRIVATE KEY")
	if err != nil {
		return nil, err
	}

	return x509.ParsePKCS1PrivateKey(d)
}
