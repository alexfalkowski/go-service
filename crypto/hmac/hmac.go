package hmac

import (
	"crypto/hmac"
	"crypto/sha512"

	"github.com/alexfalkowski/go-service/v2/crypto/errors"
	"github.com/alexfalkowski/go-service/v2/os"
)

// NewSigner constructs an HMAC-SHA-512 Signer when configuration is enabled.
//
// Disabled behavior: if cfg is nil (disabled), NewSigner returns (nil, nil).
//
// Enabled behavior: NewSigner resolves the HMAC key material via cfg.GetKey(fs) and returns a Signer
// configured to use HMAC-SHA-512 for signing and verification.
//
// Any error returned by cfg.GetKey is returned to the caller.
func NewSigner(fs *os.FS, cfg *Config) (*Signer, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}

	k, err := cfg.GetKey(fs)
	return &Signer{key: k}, err
}

// Signer signs and verifies messages using HMAC-SHA-512.
//
// The key is treated as a secret. Do not log it or expose it in errors.
type Signer struct {
	key []byte
}

// Sign computes the HMAC-SHA-512 of msg and returns the resulting MAC.
//
// This method does not fail for a valid signer configuration; it returns a nil error for API
// compatibility with other signers in this module.
func (s *Signer) Sign(msg []byte) ([]byte, error) {
	mac := hmac.New(sha512.New, s.key)
	mac.Write(msg)

	return mac.Sum(nil), nil
}

// Verify checks that sig is a valid HMAC-SHA-512 for msg.
//
// Verification uses crypto/hmac.Equal to compare MACs in constant time.
//
// It returns crypto/errors.ErrInvalidMatch when verification fails.
func (s *Signer) Verify(sig, msg []byte) error {
	mac := hmac.New(sha512.New, s.key)
	mac.Write(msg)

	if !hmac.Equal(sig, mac.Sum(nil)) {
		return errors.ErrInvalidMatch
	}

	return nil
}
