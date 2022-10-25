package test

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

// NewGenerator for test.
func NewGenerator(token string, err error) *Generator {
	return &Generator{token: token, err: err}
}

// Generator for test.
type Generator struct {
	token string
	err   error
}

func (g *Generator) Generate(ctx context.Context) ([]byte, error) {
	return []byte(g.token), g.err
}

// NewVerifier for test.
func NewVerifier(token string) *Verifier {
	return &Verifier{token: token}
}

// Verifier for test.
type Verifier struct {
	token string
}

func (v *Verifier) Verify(ctx context.Context, token []byte) (*jwt.Token, error) {
	if string(token) != v.token {
		return nil, errors.New("invalid token")
	}

	jwtToken := &jwt.Token{
		Claims: jwt.MapClaims{
			"azp": v.token,
		},
	}

	return jwtToken, nil
}
