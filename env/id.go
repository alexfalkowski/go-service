package env

import (
	"cmp"

	"github.com/alexfalkowski/go-service/id"
	"github.com/alexfalkowski/go-service/os"
)

// NewID for this service.
func NewID(generator id.Generator) ID {
	return ID(cmp.Or(os.GetVariable("SERVICE_ID"), generator.Generate()))
}

// ID of the service.
type ID string

// String representation of the ID.
func (id ID) String() string {
	return string(id)
}
