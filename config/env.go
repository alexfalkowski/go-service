package config

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/encoding/base64"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// NewENV constructs an env-based Decoder.
//
// location is the name of the environment variable to read. The environment variable value is expected to be in the form:
// "<extension>:<base64-content>" (for example "yaml:...").
//
// The extension is used to select an encoder/decoder from enc. The content must be base64-encoded.
func NewENV(location string, enc *encoding.Map) *ENV {
	kind, data := strings.CutColon(os.Getenv(location))
	return &ENV{kind: kind, data: data, enc: enc}
}

// ENV decodes configuration from an environment variable.
type ENV struct {
	enc  *encoding.Map
	kind string
	data string
}

// Decode decodes the configuration into v.
//
// It returns ErrEnvMissing if the environment variable is missing or does not contain both kind and data.
// It returns ErrNoEncoder if there is no encoder registered for the decoded kind.
// It returns an error if the base64 content cannot be decoded or if decoding into v fails.
func (e *ENV) Decode(v any) error {
	if strings.IsEmpty(e.kind) || strings.IsEmpty(e.data) {
		return ErrEnvMissing
	}

	data, err := base64.Decode(e.data)
	if err != nil {
		return err
	}

	enc := e.enc.Get(e.kind)
	if enc == nil {
		return ErrNoEncoder
	}

	return enc.Decode(bytes.NewBuffer(data), v)
}
