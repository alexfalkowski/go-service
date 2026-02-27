package nanoid

import (
	"github.com/alexfalkowski/go-service/v2/runtime"
	nanoid "github.com/matoous/go-nanoid"
)

// NewGenerator constructs a NanoID generator.
//
// The returned generator produces NanoID identifiers via github.com/matoous/go-nanoid.
// NanoIDs are URL-friendly, relatively short, random identifiers.
//
// Note: this constructor does not accept any injected randomness source; the underlying nanoid
// implementation manages its own randomness.
func NewGenerator() *Generator {
	return &Generator{}
}

// Generator generates NanoID identifiers.
type Generator struct{}

// Generate returns a newly generated NanoID string.
//
// It calls nanoid.Nanoid and returns the generated identifier.
// If NanoID generation fails, this method panics via runtime.Must.
func (n *Generator) Generate() string {
	id, err := nanoid.Nanoid()
	runtime.Must(err)
	return id
}
