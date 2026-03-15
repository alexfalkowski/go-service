package ssh

import (
	"crypto/ed25519"
	"fmt"

	"github.com/alexfalkowski/go-service/v2/crypto/errors"
	"github.com/alexfalkowski/go-service/v2/os"
	"golang.org/x/crypto/ssh"
)

// Config configures SSH key loading for Ed25519 keys used by this package.
//
// Public and Private are "source strings" resolved via os.FS.ReadSource (for example "env:NAME", "file:/path",
// or a literal key value).
//
// Expected key formats:
//   - Public: SSH authorized_keys format (parsed via ssh.ParseAuthorizedKey).
//   - Private: SSH private key format (parsed via ssh.ParseRawPrivateKey).
//
// If the provided key material is a valid SSH key but not an Ed25519 key,
// PublicKey and PrivateKey return crypto/errors.ErrInvalidKeyType.
type Config struct {
	// Public is a "source string" for the SSH public key in authorized_keys format.
	//
	// The value is resolved via os.FS.ReadSource and parsed via ssh.ParseAuthorizedKey.
	Public string `yaml:"public,omitempty" json:"public,omitempty" toml:"public,omitempty"`

	// Private is a "source string" for the SSH private key.
	//
	// The value is resolved via os.FS.ReadSource and parsed via ssh.ParseRawPrivateKey.
	Private string `yaml:"private,omitempty" json:"private,omitempty" toml:"private,omitempty"`
}

// IsEnabled reports whether SSH configuration is enabled.
//
// By convention, a nil *Config is treated as "SSH disabled" by wiring that depends on this configuration.
func (c *Config) IsEnabled() bool {
	return c != nil
}

// PublicKey resolves and parses the configured Ed25519 public key.
//
// It reads the public key data via os.FS.ReadSource and parses it as an SSH authorized key.
//
// If the parsed SSH public key is not an Ed25519 key, PublicKey returns
// crypto/errors.ErrInvalidKeyType. This can happen if the input is a valid
// authorized_keys entry but contains a different key type (for example RSA).
//
// The returned error wraps crypto/errors.ErrInvalidKeyType, so callers can use
// errors.Is to distinguish this case from read or SSH parsing errors.
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

	key, ok := parsed.(ssh.CryptoPublicKey)
	if !ok {
		return nil, fmt.Errorf("ssh: invalid public key type %T: %w", parsed, errors.ErrInvalidKeyType)
	}

	publicKey, ok := key.CryptoPublicKey().(ed25519.PublicKey)
	if !ok {
		return nil, fmt.Errorf("ssh: invalid public key type %T: %w", key.CryptoPublicKey(), errors.ErrInvalidKeyType)
	}

	return publicKey, nil
}

// PrivateKey resolves and parses the configured Ed25519 private key.
//
// It reads the private key data via os.FS.ReadSource and parses it as an SSH private key.
//
// If the parsed SSH private key is not an Ed25519 key, PrivateKey returns
// crypto/errors.ErrInvalidKeyType. This can happen if the input is a valid SSH
// private key but contains a different key type.
//
// The returned error wraps crypto/errors.ErrInvalidKeyType, so callers can use
// errors.Is to distinguish this case from read or SSH parsing errors.
func (c *Config) PrivateKey(fs *os.FS) (ed25519.PrivateKey, error) {
	data, err := fs.ReadSource(c.Private)
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParseRawPrivateKey(data)
	if err != nil {
		return nil, err
	}

	k, ok := key.(*ed25519.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("ssh: invalid private key type %T: %w", key, errors.ErrInvalidKeyType)
	}

	return *k, nil
}
