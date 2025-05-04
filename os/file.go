package os

import (
	"bytes"
	"os"
	"path/filepath"
)

// ReadFile for the name provided.
func ReadFile(name string) ([]byte, error) {
	name = ExpandPath(filepath.Clean(name))
	b, err := os.ReadFile(filepath.Clean(name))

	return bytes.TrimSpace(b), err
}

// WriteFile writes data to name with perm.
func WriteFile(name string, data []byte, perm FileMode) error {
	data = bytes.TrimSpace(data)
	name = ExpandPath(filepath.Clean(name))

	return os.WriteFile(name, data, perm)
}
