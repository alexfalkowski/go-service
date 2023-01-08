package cmd

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/alexfalkowski/go-service/file"
)

// ENV for cmd.
type ENV struct {
	location string
}

// NewENV for cmd.
func NewENV(location string) *ENV {
	return &ENV{location: location}
}

// Read for env.
func (e *ENV) Read() ([]byte, error) {
	return os.ReadFile(e.name())
}

// Write for env.
func (e *ENV) Write(data []byte, mode fs.FileMode) error {
	return os.WriteFile(e.name(), data, mode)
}

// Write for env.
func (e *ENV) Kind() string {
	return file.Extension(e.name())
}

func (e *ENV) name() string {
	return filepath.Clean(os.Getenv(e.location))
}
