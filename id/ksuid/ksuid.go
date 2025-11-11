package ksuid

import (
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/segmentio/ksuid"
)

// NewGenerator creates a new KSUID generator.
func NewGenerator() *Generator {
	return &Generator{}
}

// Generator for KSUID.
type Generator struct{}

// Generate a KSUID.
func (k *Generator) Generate() string {
	id, err := ksuid.NewRandom()
	runtime.Must(err)
	return id.String()
}
