package os

import (
	"os"
	"path/filepath"
)

// Args is an alias for os.Args.
var Args = os.Args

// ExecutableName of the running application.
func ExecutableName() string {
	path, _ := os.Executable()

	return filepath.Base(path)
}
