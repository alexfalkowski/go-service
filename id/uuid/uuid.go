package uuid

import (
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/google/uuid"
)

// NewGenerator creates a new UUID generator.
func NewGenerator() *Generator {
	return &Generator{}
}

// Generator for UUIDs.
type Generator struct{}

// Generate a UUID.
func (g *Generator) Generate() string {
	id, err := uuid.NewV7()
	runtime.Must(err)
	return id.String()
}
