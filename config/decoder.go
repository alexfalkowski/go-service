package config

import (
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/flag"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/strings"
	"go.uber.org/fx"
)

// DecoderParams for config.
type DecoderParams struct {
	fx.In

	Flags   *flag.FlagSet
	Encoder *encoding.Map
	FS      *os.FS
	Name    env.Name
}

// NewDecoder for config.
func NewDecoder(params DecoderParams) Decoder {
	kind, location := strings.CutColon(params.Flags.GetInput())

	switch kind {
	case "file":
		return NewFile(location, params.Encoder, params.FS)
	case "env":
		return NewENV(location, params.Encoder)
	default:
		return NewCommon(params.Name, params.Encoder, params.FS)
	}
}

// Decoder for config.
type Decoder interface {
	Decode(v any) error
}
