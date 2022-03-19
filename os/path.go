package os

import (
	"os"
	"path/filepath"
)

// ExecutableName of the running application.
func ExecutableName() (string, error) {
	path, err := os.Executable()
	if err != nil {
		return "", err
	}

	return filepath.Base(path), nil
}
