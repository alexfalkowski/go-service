package ed25519

import (
	"crypto/ed25519"

	crypto "github.com/alexfalkowski/go-service/v2/crypto/errors"
	"github.com/alexfalkowski/go-service/v2/crypto/pem"
)

// NewSigner for ed25519.
func NewVerifier(decoder *pem.Decoder, cfg *Config) (*Verifier, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}

	pub, err := cfg.PublicKey(decoder)
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
