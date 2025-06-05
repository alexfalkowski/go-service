package token

// NewGenerator for token.
func NewGenerator(token *Token) Generator {
	if token != nil {
		return token
	}

	return nil
}

// Generator allows the implementation of different types generators.
type Generator interface {
	// Generate a new token or error.
	Generate(aud, sub string) ([]byte, error)
}
