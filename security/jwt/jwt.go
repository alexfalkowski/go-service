package jwt

import (
	"context"

	"github.com/golang-jwt/jwt/v4"
)

// Generator allows the implementation of different types generators.
type Generator interface {
	// Generate a new token or error.
	Generate(ctx context.Context) ([]byte, error)
}

// Verifier allows the implementation of different types of verifiers.
type Verifier interface {
	// Verify a token or error.
	Verify(ctx context.Context, token []byte) (*jwt.Token, *jwt.RegisteredClaims, error)
}
