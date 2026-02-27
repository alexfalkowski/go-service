package ksuid

import (
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/segmentio/ksuid"
)

// NewGenerator constructs a KSUID generator.
//
// The returned generator produces KSUID identifiers (K-Sortable Unique Identifiers)
// via ksuid.NewRandom.
func NewGenerator() *Generator {
	return &Generator{}
}

// Generator generates KSUID identifiers.
type Generator struct{}

// Generate returns a newly generated KSUID string.
//
// It calls ksuid.NewRandom and returns the canonical string representation of the KSUID.
// If KSUID generation fails, this method panics via runtime.Must.
func (k *Generator) Generate() string {
	id, err := ksuid.NewRandom()
	runtime.Must(err)
	return id.String()
}
