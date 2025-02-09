package cmd

import (
	"io/fs"
	"os"
	"path/filepath"

	sos "github.com/alexfalkowski/go-service/os"
)

// File for cmd.
type File string

// NewFile for cmd.
func NewFile(location string) File {
	return File(location)
}

// Read for file.
func (f File) Read() ([]byte, error) {
	return os.ReadFile(f.name())
}

// Write for file.
func (f File) Write(data []byte, mode fs.FileMode) error {
	return os.WriteFile(f.name(), data, mode)
}

// Kind for file.
func (f File) Kind() string {
	return sos.PathExtension(f.name())
}

func (f File) name() string {
	return filepath.Clean(string(f))
}
