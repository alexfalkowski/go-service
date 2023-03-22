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

func (g *Generator) Generate(_ context.Context) ([]byte, error) {
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

func (v *Verifier) Verify(_ context.Context, token []byte) (*jwt.Token, *jwt.RegisteredClaims, error) {
	if string(token) != v.token {
		return nil, nil, errors.New("invalid token")
	}

	claims := &jwt.RegisteredClaims{
		Issuer:   "test",
		Subject:  "test",
		Audience: jwt.ClaimStrings{"test"},
	}

	return &jwt.Token{Claims: claims}, claims, nil
}
