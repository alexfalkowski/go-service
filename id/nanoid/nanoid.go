package nanoid

import (
	"github.com/alexfalkowski/go-service/v2/runtime"
	nanoid "github.com/matoous/go-nanoid"
)

// NewGenerator creates a new NanoID generator.
func NewGenerator() *Generator {
	return &Generator{}
}

// Generator for NanoID.
type Generator struct{}

// Generate a NanoID.
func (n *Generator) Generate() string {
	id, err := nanoid.Nanoid()
	runtime.Must(err)
	return id
}
