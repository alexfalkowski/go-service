package argon2

import (
	"github.com/alexfalkowski/go-service/crypto/errors"
	"github.com/matthewhartstonge/argon2"
)

// Algo for argon2.
type Algo interface {
	// Generate an encoded msg.
	Generate(msg string) (string, error)

	// Compare encoded with msg.
	Compare(enc, msg string) error
}

// NewAlgo for argon2.
func NewAlgo() Algo {
	return &algo{argon: argon2.DefaultConfig()}
}

type algo struct {
	argon argon2.Config
}

func (a *algo) Generate(msg string) (string, error) {
	e, err := a.argon.HashEncoded([]byte(msg))

	return string(e), err
}

func (a *algo) Compare(enc, msg string) error {
	ok, err := argon2.VerifyEncoded([]byte(msg), []byte(enc))
	if err != nil {
		return err
	}

	if !ok {
		return errors.ErrMismatch
	}

	return nil
}
