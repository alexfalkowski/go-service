package hmac

import (
	"crypto/hmac"
	"crypto/sha512"

	"github.com/alexfalkowski/go-service/v2/crypto/errors"
	"github.com/alexfalkowski/go-service/v2/os"
)

// NewSigner constructs a Signer when configuration is enabled.
//
// If cfg is disabled, it returns (nil, nil). When enabled, it loads the key material via cfg.GetKey.
// The signer uses HMAC-SHA-512 for signing and verification.
func NewSigner(fs *os.FS, cfg *Config) (*Signer, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}

	k, err := cfg.GetKey(fs)
	return &Signer{key: k}, err
}

// Signer signs and verifies messages using HMAC-SHA-512.
type Signer struct {
	key []byte
}

// Sign computes the HMAC-SHA-512 of msg.
//
// This method returns a nil error for API compatibility.
func (s *Signer) Sign(msg []byte) ([]byte, error) {
	mac := hmac.New(sha512.New, s.key)
	mac.Write(msg)

	return mac.Sum(nil), nil
}

// Verify checks that sig is a valid HMAC-SHA-512 for msg.
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
