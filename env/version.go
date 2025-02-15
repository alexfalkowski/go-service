package env

import (
	"runtime/debug"
)

// NewVersion returns the version of the application.
func NewVersion() Version {
	info, _ := debug.ReadBuildInfo()

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
