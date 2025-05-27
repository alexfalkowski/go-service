package io

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/encoding/base64"
	"github.com/alexfalkowski/go-service/v2/io"
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

// Reader from the environment variable.
func (e *ENV) Reader() io.ReadCloser {
	data, err := base64.Decode(e.data)
	if err != nil {
		return io.NewErrReadCloser(err)
	}

	return io.NopCloser(bytes.NewBuffer(data))
}

// Kind for env, which is the in the environment variable.
func (e *ENV) Kind() string {
	return e.kind
}
