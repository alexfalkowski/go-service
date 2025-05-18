package runtime

import "runtime/debug"

// Version of the binary.
func Version() string {
	info, _ := debug.ReadBuildInfo()
	if info != nil {
		return info.Main.Version
	}

	return "development"
}
