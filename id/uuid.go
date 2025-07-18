package id

import (
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/google/uuid"
)

// UUID generator.
type UUID struct{}

// Generate a UUID.
func (g *UUID) Generate() string {
	id, err := uuid.NewV7()
	runtime.Must(err)
	return id.String()
}
