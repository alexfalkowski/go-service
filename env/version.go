package env

import (
	"cmp"

	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// NewVersion returns the service version.
//
// It prefers the SERVICE_VERSION environment variable when set; otherwise it falls back to
// runtime.Version().
//
// The returned Version is not normalized at construction time; normalization (currently stripping a
// leading "v") is applied when calling Version.String.
func NewVersion() Version {
	return Version(cmp.Or(os.Getenv("SERVICE_VERSION"), runtime.Version()))
}

// Version is the service version.
//
// It is typically derived from SERVICE_VERSION (if set) or from runtime/build metadata.
type Version string

// String returns the version as a string.
//
// Normalization: if the version is prefixed with "v" (for example "v1.2.3"), the prefix is stripped
// so callers get "1.2.3". Other values are returned unchanged.
func (v Version) String() string {
	s := string(v)
	if strings.IsEmpty(s) || s[0] != 'v' {
		return s
	}

	return s[1:]
}
