package uuid

import (
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/google/uuid"
)

func init() {
	// UUIDv7 generation is on the request metadata hot path. The google/uuid
	// random pool cuts it from 2 allocs/op to 1 alloc/op in BenchmarkGenerators.
	uuid.EnableRandPool()
}

// NewGenerator constructs a UUID generator.
//
// The returned generator produces UUIDv7 identifiers (time-ordered UUIDs) via uuid.NewV7.
func NewGenerator() *Generator {
	return &Generator{}
}

// Generator generates UUID identifiers.
//
// This generator currently produces UUIDv7 values (time-ordered UUIDs).
type Generator struct{}

// Generate returns a newly generated UUIDv7 string.
//
// It calls uuid.NewV7 and returns the canonical string representation of the UUID.
// If UUID generation fails, this method panics via runtime.Must.
func (g *Generator) Generate() string {
	id, err := uuid.NewV7()
	runtime.Must(err)
	return id.String()
}
