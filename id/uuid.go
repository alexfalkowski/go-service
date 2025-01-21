package id

import (
	"github.com/google/uuid"
)

// UUID generator.
type UUID struct{}

// Generate a UUID.
func (g *UUID) Generate() string {
	return uuid.NewString()
}
