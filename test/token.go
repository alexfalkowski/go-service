package test

import (
	"errors"

	"github.com/alexfalkowski/go-service/pkg/security/token"
)

// NewGenerator for test.
func NewGenerator(token string) token.Generator {
	return &generator{token: token}
}

type generator struct {
	token string
}

func (g *generator) Generate() ([]byte, error) {
	return []byte(g.token), nil
}

// NewVerifier for test.
func NewVerifier(token string) token.Verifier {
	return &verifier{token: token}
}

type verifier struct {
	token string
}

func (v *verifier) Verify(token []byte) error {
	if string(token) != v.token {
		return errors.New("invalid token")
	}

	return nil
}
