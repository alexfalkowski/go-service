package ssh

import (
	"crypto/ed25519"

	crypto "github.com/alexfalkowski/go-service/v2/crypto/errors"
	"github.com/alexfalkowski/go-service/v2/os"
)

// NewVerifier constructs an SSH Verifier when configuration is enabled.
//
// Disabled behavior: if cfg is nil (disabled), NewVerifier returns (nil, nil).
//
// Enabled behavior: NewVerifier resolves and parses the Ed25519 public key via cfg.PublicKey(fs) and returns
// a Verifier that can validate Ed25519 signatures.
//
// Any error encountered while resolving/reading or parsing the key material is returned.
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
