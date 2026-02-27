package env

import (
	"cmp"

	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/os"
)

// NewID returns a service instance identifier.
//
// It prefers the SERVICE_ID environment variable when set (non-empty); otherwise it falls back to a
// newly generated id produced by generator.Generate().
//
// This is commonly used to distinguish service instances in logs/metrics/traces when multiple replicas
// are running.
func NewID(generator id.Generator) ID {
	return ID(cmp.Or(os.Getenv("SERVICE_ID"), generator.Generate()))
}

// ID is the service instance identifier.
type ID string

// String returns the id value as a string.
func (id ID) String() string {
	return string(id)
}
