package cmd

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/alexfalkowski/go-service/os"
)

// ErrLocationMissing for cmd.
var ErrLocationMissing = errors.New("location is missing")

// NewENV for cmd.
func NewENV(location string, fs os.FileSystem) *ENV {
	return &ENV{location: location, fs: fs}
}

// ENV for cmd.
type ENV struct {
	fs       os.FileSystem
	location string
}

// Read for env.
func (e *ENV) Read() (string, error) {
	if e.isMem() {
		_, e := e.split()
		d, err := base64.StdEncoding.DecodeString(os.GetVariable(e))

		return string(d), err
	}

	if e.name() == "" {
		return "", e.missingLocationError()
	}

	return e.fs.ReadFile(e.name())
}

// Write for env.
func (e *ENV) Write(data string, mode os.FileMode) error {
	if e.isMem() {
		_, e := e.split()

		return os.SetVariable(e, base64.StdEncoding.EncodeToString([]byte(data)))
	}

	if e.name() == "" {
		return e.missingLocationError()
	}

	return e.fs.WriteFile(e.name(), data, mode)
}

// Kind for env.
func (e *ENV) Kind() string {
	if e.isMem() {
		k, _ := e.split()

		return k
	}

	return os.PathExtension(e.name())
}

func (e *ENV) name() string {
	return os.GetVariable(e.location)
}

func (e *ENV) isMem() bool {
	return strings.Contains(e.name(), ":")
}

func (e *ENV) split() (string, string) {
	kind, env, _ := strings.Cut(e.name(), ":")

	return kind, env
}

func (e *ENV) missingLocationError() error {
	return fmt.Errorf("%s: %w", e.location, ErrLocationMissing)
}
