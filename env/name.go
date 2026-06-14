package env

import (
	"cmp"

	"github.com/alexfalkowski/go-service/v2/os"
)

// NewName returns the service name.
//
// It prefers the SERVICE_NAME environment variable when non-empty; otherwise it falls back to the
// executable name as determined by [fs.ExecutableName].
//
// The filesystem dependency exists to support consistent name derivation across environments and to
// enable tests to control the executable name behavior. Direct callers must pass a non-nil filesystem
// when SERVICE_NAME is unset or empty.
func NewName(fs *os.FS) Name {
	return Name(cmp.Or(os.Getenv("SERVICE_NAME"), fs.ExecutableName()))
}

// Name is the service name.
type Name string

// String returns the name as a string.
func (n Name) String() string {
	return string(n)
}
