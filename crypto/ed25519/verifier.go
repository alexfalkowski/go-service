package ed25519

import (
	"crypto/ed25519"

	crypto "github.com/alexfalkowski/go-service/v2/crypto/errors"
	"github.com/alexfalkowski/go-service/v2/crypto/pem"
)

// NewVerifier constructs a Verifier when configuration is enabled.
//
// If cfg is disabled, it returns (nil, nil). When enabled, it loads the public key using cfg.PublicKey.
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

// Verifier holds an Ed25519 public key used for signature verification.
type Verifier struct {
	PublicKey ed25519.PublicKey
}

// Verify verifies that sig is a valid Ed25519 signature for msg.
//
// It returns crypto.ErrInvalidMatch when verification fails.
func (v *Verifier) Verify(sig, msg []byte) error {
	ok := ed25519.Verify(v.PublicKey, msg, sig)
	if !ok {
		return crypto.ErrInvalidMatch
	}

	return nil
}
