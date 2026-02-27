package ed25519

import (
	"crypto/ed25519"

	crypto "github.com/alexfalkowski/go-service/v2/crypto/errors"
	"github.com/alexfalkowski/go-service/v2/crypto/pem"
)

// NewVerifier constructs an Ed25519 Verifier when configuration is enabled.
//
// Disabled behavior: if cfg is nil (disabled), NewVerifier returns (nil, nil).
//
// Enabled behavior: NewVerifier loads and parses the Ed25519 public key using cfg.PublicKey.
// It returns any error encountered during PEM decoding or key parsing.
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
	// PublicKey is the Ed25519 public key used by Verify.
	PublicKey ed25519.PublicKey
}

// Verify verifies that sig is a valid Ed25519 signature for msg.
//
// It returns crypto.ErrInvalidMatch when verification fails.
//
// Note: Ed25519 verification is a boolean check and does not return an error on failure. This method
// returns a sentinel error to provide a uniform verification API across crypto implementations.
func (v *Verifier) Verify(sig, msg []byte) error {
	ok := ed25519.Verify(v.PublicKey, msg, sig)
	if !ok {
		return crypto.ErrInvalidMatch
	}

	return nil
}
