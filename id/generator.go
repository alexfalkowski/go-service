package id

// Generator to generate an identifier.
type Generator interface {
	// Generate an identifier.
	Generate() string
}
