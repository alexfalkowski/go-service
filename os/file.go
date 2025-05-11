package os

import (
	"os"

	"github.com/alexfalkowski/go-service/bytes"
)

// ReadFile for the name provided.
func ReadFile(name string) ([]byte, error) {
	b, err := os.ReadFile(CleanPath(name))

	return bytes.TrimSpace(b), err
}

// WriteFile writes data to name with perm.
func WriteFile(name string, data []byte, perm FileMode) error {
	return os.WriteFile(CleanPath(name), bytes.TrimSpace(data), perm)
}
