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
func (a *Signer) Sign(msg string) (string, error) {
	e, err := a.argon.HashEncoded([]byte(msg))

	return string(e), err
}

// Verify for argon2.
func (a *Signer) Verify(sig, msg string) error {
	ok, err := argon2.VerifyEncoded([]byte(msg), []byte(sig))
	if err != nil {
		return err
	}

	if !ok {
		return errors.ErrInvalidMatch
	}

	return nil
}
