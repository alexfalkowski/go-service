package cmd

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexfalkowski/go-service/file"
)

// ErrLocationMissing for cmd.
var ErrLocationMissing = errors.New("location is missing")

// ENV for cmd.
type ENV string

// NewENV for cmd.
func NewENV(location string) ENV {
	return ENV(location)
}

// Read for env.
func (e ENV) Read() ([]byte, error) {
	if e.isMem() {
		_, e := e.split()

		return base64.StdEncoding.DecodeString(os.Getenv(e))
	}

	if e.name() == "" {
		return nil, e.missingLocationError()
	}

	return os.ReadFile(e.path())
}

// Write for env.
func (e ENV) Write(data []byte, mode fs.FileMode) error {
	if e.isMem() {
		_, e := e.split()

		return os.Setenv(e, base64.StdEncoding.EncodeToString(data))
	}

	if e.name() == "" {
		return e.missingLocationError()
	}

	return os.WriteFile(e.path(), data, mode)
}

// Kind for env.
func (e ENV) Kind() string {
	if e.isMem() {
		k, _ := e.split()

		return k
	}

	return file.Extension(e.path())
}

func (e ENV) path() string {
	return filepath.Clean(e.name())
}

func (e ENV) name() string {
	return os.Getenv(string(e))
}

func (e ENV) isMem() bool {
	return strings.Contains(e.name(), ":")
}

func (e ENV) split() (string, string) {
	s := strings.Split(e.name(), ":")

	return s[0], s[1]
}

func (e ENV) missingLocationError() error {
	return fmt.Errorf("%s: %w", string(e), ErrLocationMissing)
}
