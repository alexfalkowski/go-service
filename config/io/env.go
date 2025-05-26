package io

import (
	"github.com/alexfalkowski/go-service/v2/encoding/base64"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// NewENV for io.
func NewENV(location string) *ENV {
	kind, data := strings.CutColon(os.Getenv(location))

	return &ENV{kind: kind, data: data}
}

// ENV for io.
type ENV struct {
	kind, data string
}

// Read the data from the environment variable.
func (e *ENV) Read() ([]byte, error) {
	return base64.Decode(e.data)
}

// Kind for env, which is the in the environment variable.
func (e *ENV) Kind() string {
	return e.kind
}
