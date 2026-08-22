package config

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/encoding/base64"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/strings"
)

type envDecoder struct {
	enc      *encoding.Map
	location string
	kind     string
	data     string
}

func (e *envDecoder) Decode(v any) error {
	if strings.IsEmpty(e.kind) || strings.IsEmpty(e.data) {
		return errors.Prefix("env "+e.location, ErrEnvMissing)
	}

	data, err := base64.Decode(e.data)
	if err != nil {
		return err
	}

	enc := e.enc.Get(e.kind)
	if enc == nil {
		return errors.Prefix(strings.Join(strings.Space, "env", e.location, "kind", e.kind), ErrNoEncoder)
	}

	return enc.Decode(bytes.NewBuffer(data), v)
}
