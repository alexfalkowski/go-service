package token

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
