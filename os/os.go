package os

import (
	"os"
	"path/filepath"
)

var (
	// Args is an alias for os.Args.
	Args = os.Args

	// Stdout is an alias for os.Stdout.
	Stdout = os.Stdout
)

// ExecutableName of the running application.
func ExecutableName() string {
	return filepath.Base(Executable())
}

// ExecutableDir of the running application.
func ExecutableDir() string {
	return filepath.Dir(Executable())
}

// Executable of the running application.
func Executable() string {
	path, _ := os.Executable()

	return path
}

// UserHomeDir of the current user.
func UserHomeDir() string {
	home, _ := os.UserHomeDir()

	return home
}
