package config

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/encoding/base64"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// NewENV for config.
func NewENV(location string, enc *encoding.Map) *ENV {
	kind, data := strings.CutColon(os.Getenv(location))
	return &ENV{kind: kind, data: data, enc: enc}
}

// ENV for config.
type ENV struct {
	enc  *encoding.Map
	kind string
	data string
}

// Decode to v.
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
