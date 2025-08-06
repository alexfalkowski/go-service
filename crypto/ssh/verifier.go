package ssh

import (
	"crypto/ed25519"

	crypto "github.com/alexfalkowski/go-service/v2/crypto/errors"
	"github.com/alexfalkowski/go-service/v2/os"
)

// NewVerifier for ssh.
func NewVerifier(fs *os.FS, cfg *Config) (*Verifier, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}

	pub, err := cfg.PublicKey(fs)
	if err != nil {
		return nil, err
	}

	return &Verifier{PublicKey: pub}, nil
}

// Verifier for ssh.
type Verifier struct {
	PublicKey ed25519.PublicKey
}

// Verify for ssh.
func (v *Verifier) Verify(sig, msg []byte) error {
	ok := ed25519.Verify(v.PublicKey, msg, sig)
	if !ok {
		return crypto.ErrInvalidMatch
	}

	return nil
}
