package token

// Verifier allows the implementation of different types of verifiers.
type Verifier interface {
	// Verify a token or error.
	Verify(token []byte, aud string) (string, error)
}
