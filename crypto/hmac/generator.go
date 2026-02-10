package hmac

import (
	"github.com/alexfalkowski/go-service/v2/crypto/rand"
	"github.com/alexfalkowski/go-service/v2/errors"
)

// NewGenerator constructs a Generator that produces HMAC secrets using generator as the randomness source.
func NewGenerator(generator *rand.Generator) *Generator {
	return &Generator{generator: generator}
}

// Generator generates secrets suitable for HMAC key material.
type Generator struct {
	generator *rand.Generator
}

// Generate returns a new secret string.
//
// The returned string is 32 characters long and is drawn from the crypto/rand generator's letter set.
func (g *Generator) Generate() (string, error) {
	t, err := g.generator.GenerateText(32)
	return t, g.prefix(err)
}

func (g *Generator) prefix(err error) error {
	return errors.Prefix("hmac", err)
}
