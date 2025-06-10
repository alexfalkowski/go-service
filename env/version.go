package env

import (
	"cmp"

	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// NewVersion returns the version of the application.
func NewVersion() Version {
	return Version(cmp.Or(os.Getenv("SERVICE_VERSION"), runtime.Version()))
}

// Version of the application.
type Version string

// String of version.
func (v Version) String() string {
	s := string(v)
	if strings.IsEmpty(s) || s[0] != 'v' {
		return s
	}

	return s[1:]
}
