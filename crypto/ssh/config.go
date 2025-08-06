package ssh

import (
	"crypto/ed25519"

	"github.com/alexfalkowski/go-service/v2/os"
	"golang.org/x/crypto/ssh"
)

// Config for ssh.
type Config struct {
	Public  string `yaml:"public,omitempty" json:"public,omitempty" toml:"public,omitempty"`
	Private string `yaml:"private,omitempty" json:"private,omitempty" toml:"private,omitempty"`
}

// IsEnabled for ssh.
func (c *Config) IsEnabled() bool {
	return c != nil
}

// PublicKey ssh.
func (c *Config) PublicKey(fs *os.FS) (ed25519.PublicKey, error) {
	data, err := fs.ReadSource(c.Public)
	if err != nil {
		return nil, err
	}

	//nolint:dogsled
	parsed, _, _, _, err := ssh.ParseAuthorizedKey(data)
	if err != nil {
		return nil, err
	}

	key := parsed.(ssh.CryptoPublicKey)

	return key.CryptoPublicKey().(ed25519.PublicKey), nil
}

// PrivateKey ssh.
func (c *Config) PrivateKey(fs *os.FS) (ed25519.PrivateKey, error) {
	data, err := fs.ReadSource(c.Private)
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParseRawPrivateKey(data)
	if err != nil {
		return nil, err
	}

	k := key.(*ed25519.PrivateKey)

	return *k, nil
}
