package hmac

import "github.com/alexfalkowski/go-service/crypto/rand"

// NewGenerator for hmac.
func NewGenerator(generator *rand.Generator) *Generator {
	return &Generator{generator: generator}
}

// Generator for hmac.
type Generator struct {
	generator *rand.Generator
}

// Generate for hmac.
func (g *Generator) Generate() (string, error) {
	return g.generator.GenerateText(32)
}
