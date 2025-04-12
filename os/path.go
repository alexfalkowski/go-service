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
