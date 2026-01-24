package hmac

import (
	"github.com/alexfalkowski/go-service/v2/crypto/rand"
	"github.com/alexfalkowski/go-service/v2/errors"
)

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
	t, err := g.generator.GenerateText(32)
	return t, g.prefix(err)
}

func (g *Generator) prefix(err error) error {
	return errors.Prefix("hmac", err)
}
