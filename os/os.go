package os

import (
	"encoding/base64"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// GetVariable by key.
func GetVariable(key string) string {
	return os.Getenv(key)
}

// SetVariable of value by key.
func SetVariable(key, value string) error {
	return os.Setenv(key, value)
}

// UnsetVariable by key.
func UnsetVariable(key string) error {
	return os.Unsetenv(key)
}

// PathExtension of the specified path.
func PathExtension(path string) string {
	e := filepath.Ext(path)
	if e == "" {
		return e
	}

	return e[1:]
}

// IsNotExist if the error is os.ErrNotExist.
func IsNotExist(err error) bool {
	return errors.Is(err, os.ErrNotExist)
}

// PathExists for the path provided.
func PathExists(path string) bool {
	if _, err := os.Stat(path); IsNotExist(err) {
		return false
	}

	return true
}

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

// ExecutableName of the running application.
func ExecutableName() string {
	path, _ := os.Executable()

	return filepath.Base(path)
}
