package env

import (
	"cmp"

	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/os"
)

// NewID for this service.
func NewID(generator id.Generator) ID {
	return ID(cmp.Or(os.Getenv("SERVICE_ID"), generator.Generate()))
}

// ID of the service.
type ID string

// String representation of the ID.
func (id ID) String() string {
	return string(id)
}
