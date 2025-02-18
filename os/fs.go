package os

import (
	"os"
)

// ModeAppend will append to an existing file.
const ModeAppend = os.ModeAppend

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

// ReadFile for os.
func (*SystemFS) ReadFile(path string) (string, error) {
	return ReadFile(path)
}

// WriteFile for os.
func (*SystemFS) WriteFile(name string, data string, perm FileMode) error {
	return WriteFile(name, data, perm)
}

// PathExists for os.
func (*SystemFS) PathExists(name string) bool {
	return PathExists(name)
}

// IsNotExist for os.
func (*SystemFS) IsNotExist(err error) bool {
	return IsNotExist(err)
}
