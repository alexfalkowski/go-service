package test

import (
	"context"
	"errors"

	"github.com/alexfalkowski/go-service/meta"
)

// NewGenerator for test.
func NewGenerator(token string, err error) *Generator {
	return &Generator{token: token, err: err}
}

// Generator for test.
type Generator struct {
	err   error
	token string
}

func (g *Generator) Generate(ctx context.Context) (context.Context, []byte, error) {
	return ctx, []byte(g.token), g.err
}

// NewVerifier for test.
func NewVerifier(token string) *Verifier {
	return &Verifier{token: token}
}

// Verifier for test.
type Verifier struct {
	token string
}

func (v *Verifier) Verify(ctx context.Context, token []byte) (context.Context, error) {
	if string(token) != v.token {
		return ctx, errors.New("invalid token")
	}

	return WithTest(ctx, meta.String("auth")), nil
}
