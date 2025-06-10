package env

import (
	"cmp"

	"github.com/alexfalkowski/go-service/v2/os"
)

// NewName for this service.
func NewName(fs *os.FS) Name {
	return Name(cmp.Or(os.Getenv("SERVICE_NAME"), fs.ExecutableName()))
}

// Name of the service.
type Name string

// String representation of the name.
func (n Name) String() string {
	return string(n)
}
