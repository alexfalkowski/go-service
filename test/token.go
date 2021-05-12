package test

import (
	"context"
	"errors"

	"github.com/alexfalkowski/go-service/pkg/security/token"
)

// NewGenerator for test.
func NewGenerator(token string, err error) token.Generator {
	return &generator{token: token, err: err}
}

type generator struct {
	token string
	err   error
}

func (g *generator) Generate(ctx context.Context) ([]byte, error) {
	return []byte(g.token), g.err
}

// NewVerifier for test.
func NewVerifier(token string) token.Verifier {
	return &verifier{token: token}
}

type verifier struct {
	token string
}

func (v *verifier) Verify(ctx context.Context, token []byte) error {
	if string(token) != v.token {
		return errors.New("invalid token")
	}

	return nil
}
