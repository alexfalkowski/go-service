package ssh

import (
	"crypto/ed25519"

	"github.com/alexfalkowski/go-service/os"
	"golang.org/x/crypto/ssh"
)

// IsEnabled for ssh.
func IsEnabled(cfg *Config) bool {
	return cfg != nil
}

// Config for ssh.
type Config struct {
	Public  string `yaml:"public,omitempty" json:"public,omitempty" toml:"public,omitempty"`
	Private string `yaml:"private,omitempty" json:"private,omitempty" toml:"private,omitempty"`
}

// PublicKey ssh.
func (c *Config) PublicKey() (ed25519.PublicKey, error) {
	file, err := os.ReadFile(c.Public)
	if err != nil {
		return nil, err
	}

	//nolint:dogsled
	parsed, _, _, _, err := ssh.ParseAuthorizedKey([]byte(file))
	if err != nil {
		return nil, err
	}

	key := parsed.(ssh.CryptoPublicKey)

	return key.CryptoPublicKey().(ed25519.PublicKey), nil
}

// PrivateKey ssh.
func (c *Config) PrivateKey() (ed25519.PrivateKey, error) {
	file, err := os.ReadFile(c.Private)
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParseRawPrivateKey([]byte(file))
	if err != nil {
		return nil, err
	}

	k := key.(*ed25519.PrivateKey)

	return *k, nil
}
