package ssh

import (
	"crypto/ed25519"
	"os"

	"golang.org/x/crypto/ssh"
)

// IsEnabled for ssh.
func IsEnabled(cfg *Config) bool {
	return cfg != nil
}

type (
	// PublicKey for ssh.
	PublicKey string

	// PrivateKey for ssh.
	PrivateKey string

	// Config for ssh.
	Config struct {
		Public  PublicKey  `yaml:"public,omitempty" json:"public,omitempty" toml:"public,omitempty"`
		Private PrivateKey `yaml:"private,omitempty" json:"private,omitempty" toml:"private,omitempty"`
	}
)

// PublicKey ssh.
func (c *Config) PublicKey() (ed25519.PublicKey, error) {
	d, err := os.ReadFile(string(c.Public))
	if err != nil {
		return nil, err
	}

	//nolint:dogsled
	parsed, _, _, _, err := ssh.ParseAuthorizedKey(d)
	if err != nil {
		return nil, err
	}

	key := parsed.(ssh.CryptoPublicKey)

	return key.CryptoPublicKey().(ed25519.PublicKey), nil
}

// PrivateKey ssh.
func (c *Config) PrivateKey() (ed25519.PrivateKey, error) {
	d, err := os.ReadFile(string(c.Private))
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParseRawPrivateKey(d)
	if err != nil {
		return nil, err
	}

	k := key.(*ed25519.PrivateKey)

	return *k, nil
}
