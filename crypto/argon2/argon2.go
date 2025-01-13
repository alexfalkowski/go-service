package argon2

import (
	"github.com/alexfalkowski/go-service/crypto/algo"
	"github.com/alexfalkowski/go-service/crypto/errors"
	"github.com/matthewhartstonge/argon2"
)

// Signer for argon2.
type Signer interface {
	algo.Signer
}

// NewSigner for argon2.
func NewSigner() Signer {
	return &signer{argon: argon2.DefaultConfig()}
}

type signer struct {
	argon argon2.Config
}

func (a *signer) Sign(msg string) (string, error) {
	e, err := a.argon.HashEncoded([]byte(msg))

	return string(e), err
}

func (a *signer) Verify(sig, msg string) error {
	ok, err := argon2.VerifyEncoded([]byte(msg), []byte(sig))
	if err != nil {
		return err
	}

	if !ok {
		return errors.ErrInvalidMatch
	}

	return nil
}
