package argon2

import (
	"github.com/alexfalkowski/go-service/crypto/algo"
	"github.com/alexfalkowski/go-service/crypto/errors"
	"github.com/matthewhartstonge/argon2"
)

// Algo for argon2.
type Algo interface {
	algo.Signer
}

// NewAlgo for argon2.
func NewAlgo() Algo {
	return &argon2Algo{argon: argon2.DefaultConfig()}
}

type argon2Algo struct {
	argon argon2.Config
}

func (a *argon2Algo) Sign(msg string) (string, error) {
	e, err := a.argon.HashEncoded([]byte(msg))

	return string(e), err
}

func (a *argon2Algo) Verify(sig, msg string) error {
	ok, err := argon2.VerifyEncoded([]byte(msg), []byte(sig))
	if err != nil {
		return err
	}

	if !ok {
		return errors.ErrMismatch
	}

	return nil
}
