package uuid

import "uuid"

// NewGenerator constructs a UUID generator.
//
// The returned generator produces UUIDv7 identifiers (time-ordered UUIDs) via [uuid.NewV7].
func NewGenerator() *Generator {
	return &Generator{}
}

// Generator generates UUID identifiers.
//
// This generator currently produces UUIDv7 values (time-ordered UUIDs).
type Generator struct{}

// Generate returns a newly generated UUIDv7 string.
//
// It calls [uuid.NewV7] and returns the canonical string representation of the UUID.
// UUID generation cannot fail, so this method never panics.
func (g *Generator) Generate() string {
	return uuid.NewV7().String()
}
