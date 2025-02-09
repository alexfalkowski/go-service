package os

import (
	"os"
)

// FileMode is an alias to os.FileMode.
type FileMode = os.FileMode

// FileSystem borrows concepts from io/fs.
type FileSystem interface {
	// ReadFile for the path provided.
	ReadFile(path string) (string, error)

	// WriteFile writes data to name.
	WriteFile(name string, data string, perm FileMode) error

	// PathExists for the path provided.
	PathExists(path string) bool

	// IsNotExist whether the error is os.ErrNotExist.
	IsNotExist(err error) bool
}

// NewFS for os.
func NewFS() FileSystem {
	return &SystemFS{}
}

// SystemFS uses the underlying os.
type SystemFS struct{}

func (*SystemFS) ReadFile(path string) (string, error) {
	return ReadFile(path)
}

func (*SystemFS) WriteFile(name string, data string, perm FileMode) error {
	return WriteFile(name, data, perm)
}

func (*SystemFS) PathExists(name string) bool {
	return PathExists(name)
}

func (*SystemFS) IsNotExist(err error) bool {
	return IsNotExist(err)
}
