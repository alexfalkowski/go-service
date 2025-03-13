package argon2

import (
	"github.com/alexfalkowski/go-service/crypto/errors"
	"github.com/matthewhartstonge/argon2"
)

// NewSigner for argon2.
func NewSigner() *Signer {
	return &Signer{argon: argon2.DefaultConfig()}
}

// Signer for argon2.
type Signer struct {
	argon argon2.Config
}

// Sign for argon2.
func (s *Signer) Sign(msg []byte) ([]byte, error) {
	return s.argon.HashEncoded(msg)
}

// Verify for argon2.
func (s *Signer) Verify(sig, msg []byte) error {
	ok, err := argon2.VerifyEncoded(msg, sig)
	if err != nil {
		return err
	}

	if !ok {
		return errors.ErrInvalidMatch
	}

	return nil
}
