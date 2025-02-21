package cmd

import (
	"github.com/alexfalkowski/go-service/os"
)

// NewFile for cmd.
func NewFile(location string, fs os.FileSystem) *File {
	return &File{location: location, fs: fs}
}

// File for cmd.
type File struct {
	fs       os.FileSystem
	location string
}

// Read for file.
func (f *File) Read() (string, error) {
	return f.fs.ReadFile(f.location)
}

// Write for file.
func (f *File) Write(data string, mode os.FileMode) error {
	return f.fs.WriteFile(f.location, data, mode)
}

// Kind for file.
func (f *File) Kind() string {
	return os.PathExtension(f.location)
}
