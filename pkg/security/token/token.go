package token

import (
	"errors"
)

var (
	// ErrMissingToken for http.
	ErrMissingToken = errors.New("authorization token is not provided")
)

// Generator allows the implementation of different types generators.
type Generator interface {
	// Generate a new token or error.
	Generate() ([]byte, error)
}

// Verifier allows the implementation of different types of verifiers.
type Verifier interface {
	// Verify a token or error.
	Verify(token []byte) error
}
