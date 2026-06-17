package flag

import (
	"flag"

	"github.com/alexfalkowski/go-service/v2/strings"
)

// NewFlagSet creates a new FlagSet with the given name.
//
// The returned FlagSet wraps the standard library [flag.FlagSet] and uses [flag.ContinueOnError] so callers
// can handle parse errors explicitly (instead of terminating the process).
func NewFlagSet(name string) *FlagSet {
	set := flag.NewFlagSet(name, flag.ContinueOnError)
	return &FlagSet{FlagSet: set}
}

// FlagSet represents a set of defined flags.
//
// It wraps [flag.FlagSet] and provides optional support for the conventional go-service configuration
// flag ("-config" / "-c") used by the config subsystem to route decoding.
//
// This type is intentionally minimal: you can still use the embedded *[flag.FlagSet] to define and parse
// arbitrary flags.
type FlagSet struct {
	config *string
	*flag.FlagSet
}

// AddConfig adds the conventional configuration flag ("-config" / "-c") to the flag set.
//
// The value is treated as an opaque "kind:location" string that the config package interprets.
// Common examples include:
//   - "file:/path/to/config.yaml" (decode kind inferred from file extension)
//   - "env:MY_CONFIG" (read config payload from environment variable MY_CONFIG)
//
// The config package treats unsupported explicit "kind:location" values as invalid.
//
// The provided value is used as the default. AddConfig registers both flag names against the same
// backing value, so if both aliases are supplied during parsing, the later parsed flag wins.
//
// AddConfig must be called before [FlagSet.Parse]. Like the embedded standard library [flag.FlagSet],
// it panics if either "config" or "c" has already been registered.
func (f *FlagSet) AddConfig(value string) {
	f.config = &value
	f.StringVar(f.config, "config", value, "config location (format kind:location)")
	f.StringVar(f.config, "c", value, "config location (format kind:location)")
}

// GetConfig returns the configured config flag ("-config" / "-c") value.
//
// After [FlagSet.AddConfig], GetConfig returns the default value until [FlagSet.Parse] updates it.
// GetConfig is safe to call before AddConfig; in that case it returns an empty string.
func (f *FlagSet) GetConfig() string {
	if f.config != nil {
		return *f.config
	}
	return strings.Empty
}
