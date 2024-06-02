package rsa

import (
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

// GetPublic for rsa.
func (c *Config) GetPublic() ([]byte, error) {
	return pem.Decode(string(c.Public), "RSA PUBLIC KEY")
}

// GetPrivate for rsa.
func (c *Config) GetPrivate() ([]byte, error) {
	return pem.Decode(string(c.Private), "RSA PRIVATE KEY")
}
