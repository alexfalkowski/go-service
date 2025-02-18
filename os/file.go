package os

import (
	"os"
	"path/filepath"
	"strings"
)

// ReadFile for the path provided.
func ReadFile(path string) (string, error) {
	b, err := os.ReadFile(filepath.Clean(path))

	return strings.TrimSpace(string(b)), err
}

// WriteFile writes data to name with perm.
func WriteFile(name string, data string, perm FileMode) error {
	data = strings.TrimSpace(data)

	return os.WriteFile(name, []byte(data), perm)
}

// Remove a file or empty folder.
func Remove(name string) error {
	return os.Remove(name)
}
