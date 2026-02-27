// Package config provides configuration decoding, validation, and DI wiring for go-service.
//
// This package exposes a `Decoder` abstraction and multiple decoder implementations that load
// configuration from different sources. The source is selected by the "-i" flag (see `NewDecoder` and
// `flag.FlagSet.GetInput`).
//
// # Input routing (-i flag)
//
// `NewDecoder` dispatches based on the value of "-i":
//
//   - "file:<path>": loads configuration from the file at <path>. The decoder chooses the parser based
//     on the file extension.
//   - "env:<ENV_VAR>": loads configuration from the environment variable named <ENV_VAR>. The variable
//     value must be formatted as "<extension>:<base64-content>" (for example "yaml:...").
//   - otherwise: uses the default lookup, searching for "<serviceName>.{yaml,yml,toml,json}" in common
//     locations (executable directory, user config dir, and /etc).
//
// # Decoding and validation
//
// For typed configuration, use `NewConfig[T]`, which:
//   - decodes into a newly allocated `*T` using a `Decoder`,
//   - rejects empty decoded values by returning `ErrInvalidConfig`, and
//   - validates the decoded value using `Validator` (go-playground/validator).
//
// # DI wiring
//
// `Module` wires the decoder, validator, and a standard top-level `Config` into Fx/Dig, and also
// provides constructors for commonly-used sub-config projections.
package config
