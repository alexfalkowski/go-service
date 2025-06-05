package token

// Generator allows the implementation of different types generators.
type Generator interface {
	// Generate a new token or error.
	Generate(aud, sub string) ([]byte, error)
}
