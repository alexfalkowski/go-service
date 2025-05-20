package test

import (
	"context"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/strings"
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
	return ctx, strings.Bytes(g.token), g.err
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
	if bytes.String(token) != v.token {
		return ctx, ErrInvalid
	}

	return ctx, nil
}
