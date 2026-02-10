package token

// Verifier verifies tokens.
type Verifier interface {
	// Verify validates token for the given audience and returns the subject.
	Verify(token []byte, aud string) (string, error)
}
