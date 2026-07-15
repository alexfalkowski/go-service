package env

import (
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// NewID returns a service instance identifier.
//
// It prefers the SERVICE_ID environment variable when set (non-empty); otherwise it falls back to a
// newly generated id produced by generator.Generate().
//
// Direct callers must pass a non-nil generator when SERVICE_ID is unset or empty. The standard module
// wiring supplies this generator through the id module.
//
// This is commonly used to distinguish service instances in logs/metrics/traces when multiple replicas
// are running.
func NewID(generator id.Generator) ID {
	if id := os.Getenv("SERVICE_ID"); !strings.IsEmpty(id) {
		return ID(id)
	}

	return ID(generator.Generate())
}

// ID is the service instance identifier.
type ID string

// String returns the id value as a string.
func (id ID) String() string {
	return string(id)
}
