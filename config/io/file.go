package io

import "github.com/alexfalkowski/go-service/v2/os"

// NewFile for io.
func NewFile(location string, fs *os.FS) *File {
	return &File{location: location, fs: fs}
}

// File reads and writes files.
type File struct {
	fs       *os.FS
	location string
}

// Read from a file.
func (f *File) Read() ([]byte, error) {
	return f.fs.ReadFile(f.location)
}

// Write to as file.
func (f *File) Write(data []byte, mode os.FileMode) error {
	return f.fs.WriteFile(f.location, data, mode)
}

// Kind for file, which is the file extension.
func (f *File) Kind() string {
	return f.fs.PathExtension(f.location)
}
