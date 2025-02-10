package os

import (
	"errors"
	"os"
	"path/filepath"
)

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
