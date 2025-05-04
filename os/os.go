package os

import (
	"os"
	"path/filepath"
)

var (
	// Args is an alias for os.Args.
	Args = os.Args

	// Create is an alias for os.Create.
	Create = os.Create

	// MkdirAll is an alias for os.MkdirAll.
	MkdirAll = os.MkdirAll

	// Remove is an alias for os.Remove.
	Remove = os.Remove

	// RemoveAll is an alias for os.RemoveAll.
	RemoveAll = os.RemoveAll

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
