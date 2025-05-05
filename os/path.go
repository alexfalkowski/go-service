package os

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/alexfalkowski/go-service/strings"
)

// PathExtension of the specified path.
func PathExtension(path string) string {
	e := filepath.Ext(path)
	if strings.IsEmpty(e) {
		return e
	}

	return e[1:]
}

// IsNotExist if the error is os.ErrNotExist.
func IsNotExist(err error) bool {
	return errors.Is(err, os.ErrNotExist)
}

// PathExists for the name provided.
func PathExists(name string) bool {
	if _, err := os.Stat(name); IsNotExist(err) {
		return false
	}

	return true
}

// ExpandPath will append the home dir if path starts with ~.
func ExpandPath(path string) string {
	dir := UserHomeDir()
	if strings.IsEmpty(dir) || len(path) == 0 || path[0] != '~' {
		return path
	}

	return filepath.Join(dir, path[1:])
}

// CleanPath makes sure that the path is expanded and in a clean format.
func CleanPath(path string) string {
	return ExpandPath(filepath.Clean(path))
}
