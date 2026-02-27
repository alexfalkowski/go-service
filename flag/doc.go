// Package flag provides helpers for defining and parsing command-line flags in go-service.
//
// This package is a small wrapper around the standard library `flag` package that standardizes a few
// flag conventions used across go-service, especially configuration source selection.
//
// # Configuration input flag (-i)
//
// Many go-service applications accept an "-i" flag that selects where configuration should be loaded from.
// The conventional value format is:
//
//	"kind:location"
//
// Common examples:
//   - "file:/path/to/config.yaml" to read configuration from a file (decoder selected by extension)
//   - "env:MY_CONFIG" to read configuration from an environment variable (typically "<ext>:<base64-content>")
//
// The `FlagSet` type in this package supports installing this convention via `(*FlagSet).AddInput` and
// retrieving it via `(*FlagSet).GetInput`. The config subsystem (`config.NewDecoder`) consumes this value
// to route to the appropriate decoder.
//
// Start with `NewFlagSet`, `FlagSet.AddInput`, and `FlagSet.GetInput`.
package flag
