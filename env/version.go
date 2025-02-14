package env

import (
	"runtime/debug"
)

// NewVersion for env.
func NewVersion(info *debug.BuildInfo) Version {
	return Version(info.Main.Version)
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
