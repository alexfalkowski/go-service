package ssh

import (
	"crypto/ed25519"

	"github.com/alexfalkowski/go-service/v2/os"
)

// NewSigner for ssh.
func NewSigner(fs *os.FS, cfg *Config) (*Signer, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}

	pri, err := cfg.PrivateKey(fs)
	if err != nil {
		return nil, err
	}

	return &Signer{PrivateKey: pri}, nil
}

// Signer for ssh.
type Signer struct {
	PrivateKey ed25519.PrivateKey
}

// Sign for ssh.
func (s *Signer) Sign(msg []byte) ([]byte, error) {
	return ed25519.Sign(s.PrivateKey, msg), nil
}
