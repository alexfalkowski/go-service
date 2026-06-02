// Package flag provides helpers for defining and parsing command-line flags in go-service.
//
// This package is a small wrapper around the standard library [flag] package that standardizes a few
// flag conventions used across go-service, especially configuration source selection.
//
// # Configuration flag (-config / -c)
//
// Many go-service applications accept "-config" and "-c" flags that select where configuration should be loaded from.
// The conventional value format is:
//
//	"kind:location"
//
// Common examples:
//   - "file:/path/to/config.yaml" to read configuration from a file (decoder selected by extension)
//   - "env:MY_CONFIG" to read configuration from an environment variable (typically "<ext>:<base64-content>")
//
// The [FlagSet] type in this package supports installing this convention via [FlagSet.AddConfig] and
// retrieving it via [FlagSet.GetConfig]. The config subsystem ([config.NewDecoder]) consumes this value
// to route to the appropriate decoder.
//
// Start with [NewFlagSet], [FlagSet.AddConfig], and [FlagSet.GetConfig].
package flag
