package env

import (
	"cmp"

	"github.com/alexfalkowski/go-service/os"
)

// NewVersion returns the version of the application.
func NewVersion() Version {
	return Version(cmp.Or(os.GetVariable("SERVICE_VERSION"), "development"))
}

// Version of the application.
type Version string

// String of version.
func (v Version) String() string {
	s := string(v)
	if s == "" || s[0] != 'v' {
		return s
	}

	return s[1:]
}
