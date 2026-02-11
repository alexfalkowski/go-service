package ssh

import (
	"crypto/ed25519"

	"github.com/alexfalkowski/go-service/v2/os"
	"golang.org/x/crypto/ssh"
)

// Config configures SSH key loading for Ed25519 keys.
//
// Public and Private are "source strings" that are read by os.FS.ReadSource (for example "env:NAME", "file:/path",
// or a literal key value).
//
// Expected key formats:
//   - Public: SSH authorized_keys format (parsed via ssh.ParseAuthorizedKey).
//   - Private: SSH private key format (parsed via ssh.ParseRawPrivateKey).
//
// Note: Key parsing uses type assertions. If the provided key material is not an Ed25519 SSH key,
// PublicKey/PrivateKey may panic due to type assertions.
type Config struct {
	// Public is a "source string" for the SSH public key in authorized_keys format.
	//
	// It is read via os.FS.ReadSource and parsed via ssh.ParseAuthorizedKey.
	Public string `yaml:"public,omitempty" json:"public,omitempty" toml:"public,omitempty"`

	// Private is a "source string" for the SSH private key.
	//
	// It is read via os.FS.ReadSource and parsed via ssh.ParseRawPrivateKey.
	Private string `yaml:"private,omitempty" json:"private,omitempty" toml:"private,omitempty"`
}

// IsEnabled reports whether SSH configuration is enabled.
//
// A nil config is considered disabled.
func (c *Config) IsEnabled() bool {
	return c != nil
}

// PublicKey loads and parses the configured Ed25519 public key.
//
// It reads the public key data via os.FS.ReadSource and parses it as an SSH authorized key.
// If the parsed key is not an Ed25519 key, this function will panic due to type assertions.
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

// PrivateKey loads and parses the configured Ed25519 private key.
//
// It reads the private key data via os.FS.ReadSource and parses it as an SSH private key.
// If the parsed key is not an Ed25519 key, this function will panic due to type assertions.
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
