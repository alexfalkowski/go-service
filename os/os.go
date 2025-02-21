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
	path, _ := os.Executable()

	return filepath.Base(path)
}
