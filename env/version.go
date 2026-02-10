package env

import (
	"cmp"

	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// NewVersion returns the service version.
//
// It prefers the SERVICE_VERSION environment variable when set; otherwise it falls back to runtime.Version().
func NewVersion() Version {
	return Version(cmp.Or(os.Getenv("SERVICE_VERSION"), runtime.Version()))
}

// Version of the application.
type Version string

// String returns the version as a string.
//
// If the version is prefixed with "v" (for example "v1.2.3"), the prefix is stripped.
func (v Version) String() string {
	s := string(v)
	if strings.IsEmpty(s) || s[0] != 'v' {
		return s
	}

	return s[1:]
}
