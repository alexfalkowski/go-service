package hmac

import (
	"crypto/hmac"
	"crypto/sha512"

	"github.com/alexfalkowski/go-service/v2/crypto/errors"
	"github.com/alexfalkowski/go-service/v2/os"
)

// NewSigner for hmac.
func NewSigner(fs *os.FS, cfg *Config) (*Signer, error) {
	if !cfg.IsEnabled() {
		return nil, nil
	}

	k, err := cfg.GetKey(fs)
	return &Signer{key: k}, err
}

// Signer for hmac.
type Signer struct {
	key []byte
}

// Sign for hmac.
func (s *Signer) Sign(msg []byte) ([]byte, error) {
	mac := hmac.New(sha512.New, s.key)
	mac.Write(msg)

	return mac.Sum(nil), nil
}

// Verify for hmac.
func (s *Signer) Verify(sig, msg []byte) error {
	mac := hmac.New(sha512.New, s.key)
	mac.Write(msg)

	if !hmac.Equal(sig, mac.Sum(nil)) {
		return errors.ErrInvalidMatch
	}

	return nil
}
