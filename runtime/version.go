package runtime

import "runtime/debug"

// Version returns the build version of the running binary.
//
// It reads build info via [debug.ReadBuildInfo] and returns the main module
// version when build info is available. The returned value is the build-info
// value as-is, including "(devel)" for local development builds. If build info
// is unavailable, Version returns "development".
func Version() string {
	info, _ := debug.ReadBuildInfo()
	if info != nil {
		return info.Main.Version
	}
	return "development"
}
