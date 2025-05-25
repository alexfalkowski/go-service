package io

import (
	"github.com/alexfalkowski/go-service/v2/encoding/base64"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// NewENV for io.
func NewENV(location string) *ENV {
	return &ENV{location: location}
}

// ENV reads and writes from environment variables to files.
type ENV struct {
	location string
}

// Read a file from an environment variable.
// The contents of the file could be inside the environment variable.
func (e *ENV) Read() ([]byte, error) {
	_, env := e.split()

	return base64.Decode(os.Getenv(env))
}

// Kind for env, which is the file extension or defined in the environment variable.
func (e *ENV) Kind() string {
	kind, _ := e.split()

	return kind
}

func (e *ENV) name() string {
	return os.Getenv(e.location)
}

func (e *ENV) split() (string, string) {
	kind, env, _ := strings.Cut(e.name(), ":")

	return kind, env
}
