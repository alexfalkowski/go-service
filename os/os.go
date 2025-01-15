package os

import (
	"encoding/base64"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexfalkowski/go-service/runtime"
)

// FileExists for the path provided.
func FileExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}

// ReadFile for the path provided.
func ReadFile(path string) (string, error) {
	b, err := os.ReadFile(filepath.Clean(path))

	return strings.TrimSpace(string(b)), err
}

// MustReadFile is like ReadFile, though it panics on error.
func MustReadFile(path string) string {
	s, err := ReadFile(path)
	runtime.Must(err)

	return s
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

// ExecutableName of the running application.
func ExecutableName() string {
	path, _ := os.Executable()

	return filepath.Base(path)
}
