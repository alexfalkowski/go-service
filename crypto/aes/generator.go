package aes

import "github.com/alexfalkowski/go-service/v2/crypto/rand"

// NewGenerator for aes.
func NewGenerator(generator *rand.Generator) *Generator {
	return &Generator{generator: generator}
}

// Generator for aes.
type Generator struct {
	generator *rand.Generator
}

// Generate for aes.
func (g *Generator) Generate() (string, error) {
	return g.generator.GenerateText(32)
}
