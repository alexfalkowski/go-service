package aes

import "github.com/alexfalkowski/go-service/v2/crypto/rand"

// NewGenerator constructs a Generator that produces AES key material.
func NewGenerator(generator *rand.Generator) *Generator {
	return &Generator{generator: generator}
}

// Generator generates AES keys using the shared random generator.
type Generator struct {
	generator *rand.Generator
}

// Generate returns a base64 text key of length suitable for AES-256.
func (g *Generator) Generate() (string, error) {
	return g.generator.GenerateText(32)
}
