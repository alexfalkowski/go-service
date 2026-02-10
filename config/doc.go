// Package config provides configuration decoding and wiring for go-service.
//
// This package exposes a Decoder abstraction and implementations that load configuration from:
//   - file:<path> via the file decoder,
//   - env:<ENV_VAR> via the env decoder (expects "<extension>:<base64-content>"), or
//   - a default lookup that searches for "<serviceName>.{yaml,yml,toml,json}" in common locations.
//
// The input source is selected by the "-i" flag (see NewDecoder and flag.FlagSet.GetInput).
package config
