// Package config provides configuration decoding, validation, and DI wiring for go-service.
//
// This package exposes a [Decoder] abstraction and multiple decoder implementations that load
// configuration from different sources. The source is selected by the "-config" / "-c" flag (see [NewDecoder] and
// `flag.FlagSet.GetConfig`).
//
// # Config routing (-config / -c flag)
//
// [NewDecoder] dispatches based on the value of "-config" / "-c":
//
//   - "file:<path>": loads configuration from the file at <path>. The decoder chooses the parser based
//     on the file extension.
//   - "env:<ENV_VAR>": loads configuration from the environment variable named <ENV_VAR>. The variable
//     value must be formatted as "<extension>:<base64-content>" (for example "yaml:...").
//   - otherwise: uses the default lookup, searching for "<serviceName>.{yaml,yml,hjson,toml,json}" in common
//     locations (executable directory, user config dir, and /etc).
//
// Default lookup intentionally requires a valid user configuration directory environment before lookup starts.
// It may resolve the user config directory before probing file candidates, so runtimes using default lookup
// should provide HOME or XDG_CONFIG_HOME. Services that do not want that environment contract should pass an
// explicit source with "file:<path>" or "env:<ENV_VAR>".
//
// # Decoding and validation
//
// For typed configuration, use `NewConfig[T]`, which:
//   - decodes into a newly allocated `*T` using a [Decoder],
//   - rejects empty decoded values by returning [ErrInvalidConfig], and
//   - validates the decoded value using [Validator] (go-playground/validator).
//
// # DI wiring
//
// [Module] wires the decoder, validator, and a standard top-level [Config] into [go.uber.org/fx]/[go.uber.org/dig], and also
// provides constructors for commonly-used sub-config projections.
//
// In normal service applications, this package is consumed through higher-level bundles such as
// [github.com/alexfalkowski/go-service/v2/module.Server] or [github.com/alexfalkowski/go-service/v2/module.Client] from `go-service-template`, which also include the standard
// encoder registrations needed by the config decoders. Custom or partial wiring is still supported,
// but advanced compositions are responsible for registering any required encoders themselves.
package config
