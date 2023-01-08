package cmd

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/alexfalkowski/go-service/file"
)

// File for cmd.
type File struct {
	location string
}

// NewFile for cmd.
func NewFile(location string) *File {
	return &File{location: location}
}

// Read for file.
func (f *File) Read() ([]byte, error) {
	return os.ReadFile(f.name())
}

// Write for file.
func (f *File) Write(data []byte, mode fs.FileMode) error {
	return os.WriteFile(f.name(), data, mode)
}

// Write for file.
func (f *File) Kind() string {
	return file.Extension(f.name())
}

func (f *File) name() string {
	return filepath.Clean(f.location)
}
