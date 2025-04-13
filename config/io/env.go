package io

import (
	"encoding/base64"
	"fmt"

	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/strings"
)

// NewENV for io.
func NewENV(location string, fs os.FileSystem) *ENV {
	return &ENV{location: location, fs: fs}
}

// ENV reads and writes from environment variables to files.
type ENV struct {
	fs       os.FileSystem
	location string
}

// Valid checks if the location is present.
func (e *ENV) Valid() bool {
	if e.isMemory() {
		_, e := e.split()

		return os.VariableExists(e)
	}

	return e.fs.PathExists(e.name())
}

// Read a file from an environment variable.
// The contents of the file could be inside the environment variable.
func (e *ENV) Read() ([]byte, error) {
	if e.isMemory() {
		_, e := e.split()

		return base64.StdEncoding.DecodeString(os.GetVariable(e))
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

		return os.SetVariable(e, base64.StdEncoding.EncodeToString(data))
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

	return os.PathExtension(e.name())
}

func (e *ENV) name() string {
	return os.GetVariable(e.location)
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
