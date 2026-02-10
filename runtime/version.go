package runtime

import "runtime/debug"

// Version returns the build version of the running binary.
//
// It reads build info via debug.ReadBuildInfo and returns the main module version when available.
// If build info is unavailable, it returns "development".
func Version() string {
	info, _ := debug.ReadBuildInfo()
	if info != nil {
		return info.Main.Version
	}
	return "development"
}
