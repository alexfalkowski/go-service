package token

import (
	"context"
)

// NewGenerator for token.
func NewGenerator() Generator {
	return &none{}
}

// NewVerifier for token.
func NewVerifier() Verifier {
	return &none{}
}

type none struct{}

func (n *none) Generate(ctx context.Context) (context.Context, []byte, error) {
	return ctx, nil, nil
}

func (n *none) Verify(ctx context.Context, _ []byte) (context.Context, error) {
	return ctx, nil
}
