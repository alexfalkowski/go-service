package token

// Generator generates tokens.
type Generator interface {
	// Generate creates a new token for the given audience and subject.
	Generate(aud, sub string) ([]byte, error)
}
