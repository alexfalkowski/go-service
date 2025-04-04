package ed25519

import (
	"crypto/ed25519"

	crypto "github.com/alexfalkowski/go-service/crypto/errors"
)

// NewSigner for ed25519.
func NewVerifier(cfg *Config) (*Verifier, error) {
	if !IsEnabled(cfg) {
		return nil, nil
	}

	pub, err := cfg.PublicKey()
	if err != nil {
		return nil, err
	}

	return &Verifier{PublicKey: pub}, nil
}

// Verifier for ed25519.
type Verifier struct {
	PublicKey ed25519.PublicKey
}

// Verify for ed25519.
func (v *Verifier) Verify(sig, msg []byte) error {
	ok := ed25519.Verify(v.PublicKey, msg, sig)
	if !ok {
		return crypto.ErrInvalidMatch
	}

	return nil
}
