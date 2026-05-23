package ed25519

import (
	"crypto/ed25519"
	"fmt"

	"github.com/alexfalkowski/go-service/v2/crypto/errors"
)

// PrivateKeySize is the size, in bytes, of an Ed25519 private key.
const PrivateKeySize = ed25519.PrivateKeySize

// PrivateKey is the type of an Ed25519 private key.
type PrivateKey = ed25519.PrivateKey

// ValidatePrivateKey reports whether key has the Ed25519 private-key size required for signing.
func ValidatePrivateKey(key ed25519.PrivateKey) error {
	if len(key) != PrivateKeySize {
		return fmt.Errorf("ed25519: invalid private key size %d: %w", len(key), errors.ErrInvalidKeySize)
	}

	return nil
}

// Sign signs msg with key.
func Sign(key PrivateKey, msg []byte) []byte {
	return ed25519.Sign(key, msg)
}
