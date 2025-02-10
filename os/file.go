package os

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
)

// ReadFile for the path provided.
func ReadFile(path string) (string, error) {
	b, err := os.ReadFile(filepath.Clean(path))

	return strings.TrimSpace(string(b)), err
}

// ReadBase64File for the path provided.
func ReadBase64File(path string) (string, error) {
	s, err := ReadFile(path)
	if err != nil {
		return "", err
	}

	dc, err := base64.StdEncoding.DecodeString(s)

	return string(dc), err
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
