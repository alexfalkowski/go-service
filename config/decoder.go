package config

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/flag"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// DecoderParams defines the dependencies used to construct a Decoder.
//
// It is intended for dependency injection ([go.uber.org/fx]/[go.uber.org/dig]). The default wiring is provided by [Module].
type DecoderParams struct {
	di.In

	// Flags is the parsed flag set used to select the configuration source.
	// NewDecoder reads the "-config" / "-c" flag via Flags.GetConfig.
	Flags *flag.FlagSet

	// Encoder is the registry of decoders keyed by kind/extension (e.g. "yaml", "json", "hjson", "toml").
	Encoder *encoding.Map

	// FS is the filesystem used for configuration file lookup and reading.
	// It is used by the file decoder and the default lookup decoder.
	FS *os.FS

	// Name is the service name used by the default lookup decoder to locate "<serviceName>.<ext>".
	Name env.Name
}

// NewDecoder constructs a Decoder based on the configured input source.
//
// Routing is controlled by the "-config" / "-c" flag (see [flag.FlagSet.GetConfig]). The value supports a
// "kind:location" format:
//
//   - "file:<path>": uses the file decoder to read from <path>. The file extension determines which
//     encoder is used (e.g. ".yaml" -> "yaml", ".hjson" -> "hjson").
//   - "env:<ENV_VAR>": uses the env decoder to read from the environment variable <ENV_VAR>.
//     The variable value must be formatted as "<extension>:<base64-content>" (e.g. "yaml:..." or "hjson:...").
//   - unprefixed values (including empty): use the default lookup decoder, which searches common
//     locations for "<serviceName>.{yaml,yml,hjson,toml,json}".
//   - unsupported explicit "kind:location" values fail with [ErrInvalidSource].
//
// The returned Decoder is safe for repeated calls to Decode; underlying behavior depends on the
// selected implementation.
func NewDecoder(params DecoderParams) Decoder {
	kind, location, found := strings.CutColon(params.Flags.GetConfig())
	if !found {
		return NewDefault(params.Name, params.Encoder, params.FS)
	}

	switch kind {
	case "file":
		return NewFile(location, params.Encoder, params.FS)
	case "env":
		return NewENV(location, params.Encoder)
	default:
		return invalidSourceDecoder{source: params.Flags.GetConfig()}
	}
}

// Decoder loads configuration from some source and decodes it into a destination value.
//
// Implementations typically read raw configuration bytes, select a decoding strategy based on a kind
// (e.g. file extension), and unmarshal into the provided destination.
type Decoder interface {
	// Decode reads configuration from the underlying source and decodes it into v.
	//
	// The v parameter should be a pointer to the destination type. Implementations may return
	// errors for missing inputs, unknown/unsupported kinds, I/O failures, or decode/unmarshal failures.
	Decode(v any) error
}
