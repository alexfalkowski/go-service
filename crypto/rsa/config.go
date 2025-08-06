package rsa

import (
	"crypto/rsa"
	"crypto/x509"

	"github.com/alexfalkowski/go-service/v2/crypto/pem"
)

// Config for rsa.
type Config struct {
	Public  string `yaml:"public,omitempty" json:"public,omitempty" toml:"public,omitempty"`
	Private string `yaml:"private,omitempty" json:"private,omitempty" toml:"private,omitempty"`
}

// IsEnabled for rsa.
func (c *Config) IsEnabled() bool {
	return c != nil
}

// PublicKey rsa.
func (c *Config) PublicKey(decoder *pem.Decoder) (*rsa.PublicKey, error) {
	d, err := decoder.Decode(c.Public, "RSA PUBLIC KEY")
	if err != nil {
		return nil, err
	}

	return x509.ParsePKCS1PublicKey(d)
}

// PrivateKey rsa.
func (c *Config) PrivateKey(decoder *pem.Decoder) (*rsa.PrivateKey, error) {
	d, err := decoder.Decode(c.Private, "RSA PRIVATE KEY")
	if err != nil {
		return nil, err
	}

	return x509.ParsePKCS1PrivateKey(d)
}
