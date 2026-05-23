package aes

import "github.com/alexfalkowski/go-service/v2/crypto/rand"

// NewGenerator constructs a Generator that produces AES key material.
//
// The returned Generator uses the shared cryptographically-secure random generator.
func NewGenerator(generator *rand.Generator) *Generator {
	return &Generator{generator: generator}
}

// Generator generates AES keys using the shared random generator.
type Generator struct {
	generator *rand.Generator
}

// Generate returns a 32-character alphanumeric key suitable for AES-256.
//
// The generated key is safe to store in text configuration and can be used
// directly as 32-byte AES key material.
func (g *Generator) Generate() (string, error) {
	return g.generator.GenerateText(32)
}
