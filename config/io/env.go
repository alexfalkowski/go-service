package io

import (
	"fmt"

	"github.com/alexfalkowski/go-service/v2/encoding/base64"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// NewENV for io.
func NewENV(location string, fs *os.FS) *ENV {
	return &ENV{location: location, fs: fs}
}

// ENV reads and writes from environment variables to files.
type ENV struct {
	fs       *os.FS
	location string
}

// Read a file from an environment variable.
// The contents of the file could be inside the environment variable.
func (e *ENV) Read() ([]byte, error) {
	if e.isMemory() {
		_, e := e.split()

		return base64.Decode(os.Getenv(e))
	}

	if strings.IsEmpty(e.name()) {
		return nil, e.error()
	}

	return e.fs.ReadFile(e.name())
}

// Write a file from an environment variable.
// The contents of the file could be written to the environment variable.
func (e *ENV) Write(data []byte, mode os.FileMode) error {
	if e.isMemory() {
		_, e := e.split()

		return os.Setenv(e, base64.Encode(data))
	}

	if strings.IsEmpty(e.name()) {
		return e.error()
	}

	return e.fs.WriteFile(e.name(), data, mode)
}

// Kind for env, which is the file extension or defined in the environment variable.
func (e *ENV) Kind() string {
	if e.isMemory() {
		k, _ := e.split()

		return k
	}

	return e.fs.PathExtension(e.name())
}

func (e *ENV) name() string {
	return os.Getenv(e.location)
}

func (e *ENV) isMemory() bool {
	return strings.Contains(e.name(), ":")
}

func (e *ENV) split() (string, string) {
	kind, env, _ := strings.Cut(e.name(), ":")

	return kind, env
}

func (e *ENV) error() error {
	return fmt.Errorf("%s: %w", e.location, ErrLocationMissing)
}
