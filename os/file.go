package os

import (
	"bytes"
	"os"
	"path/filepath"
)

// ReadFile for the path provided.
func ReadFile(name string) ([]byte, error) {
	b, err := os.ReadFile(filepath.Clean(name))

	return bytes.TrimSpace(b), err
}

// WriteFile writes data to name with perm.
func WriteFile(name string, data []byte, perm FileMode) error {
	data = bytes.TrimSpace(data)

	return os.WriteFile(name, data, perm)
}
