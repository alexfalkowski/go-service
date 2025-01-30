package token

import (
	"context"
)

// Verifier allows the implementation of different types of verifiers.
type Verifier interface {
	// Verify a token or error.
	Verify(ctx context.Context, token []byte) (context.Context, error)
}
