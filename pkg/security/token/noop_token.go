package token

// NewTokenGenerator for security.
func NewNoopGenerator() Generator {
	return &noopToken{}
}

// NewTokenVerifier for security.
func NewNoopVerifier() Verifier {
	return &noopToken{}
}

// NoopToken satisfies Token and does nothing.
type noopToken struct {
}

func (t *noopToken) Generate() ([]byte, error) {
	return nil, nil
}

func (t *noopToken) Verify(token []byte) error {
	return nil
}
