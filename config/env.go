package config

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/encoding/base64"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// NewENV constructs an environment-variable based Decoder.
//
// location is the name of the environment variable to read.
//
// # Expected format
//
// The environment variable value must be formatted as:
//
//	"<kind>:<base64-content>"
//
// Where:
//   - kind selects the decoder from enc (for example "yaml", "yml", "toml", or "json").
//   - base64-content is the configuration content encoded using base64.
//
// Example:
//
//	CONFIG="yaml:PHlhbWw+Li4u" // base64 payload truncated for brevity
//
// NewENV reads the environment variable once at construction time and stores the parsed kind and data.
func NewENV(location string, enc *encoding.Map) *ENV {
	kind, data := strings.CutColon(os.Getenv(location))
	return &ENV{kind: kind, data: data, enc: enc}
}

// ENV decodes configuration from an environment variable.
//
// It expects the variable value to be formatted as "<kind>:<base64-content>" (see NewENV).
type ENV struct {
	enc  *encoding.Map
	kind string
	data string
}

// Decode decodes the configuration into v.
//
// The destination v should be a pointer to the target configuration type.
//
// Errors:
//   - ErrEnvMissing if the environment variable is missing or malformed (missing kind or data).
//   - ErrNoEncoder if there is no encoder registered for the decoded kind.
//   - Any base64 decode error if the content cannot be decoded.
//   - Any decode/unmarshal error returned by the selected encoder.
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
