package flag

import (
	"flag"

	"github.com/alexfalkowski/go-service/v2/strings"
)

// NewFlagSet creates a new FlagSet with the given name.
//
// The returned FlagSet wraps the standard library flag.FlagSet and uses flag.ContinueOnError so callers
// can handle parse errors explicitly (instead of terminating the process).
func NewFlagSet(name string) *FlagSet {
	set := flag.NewFlagSet(name, flag.ContinueOnError)
	return &FlagSet{FlagSet: set}
}

// FlagSet represents a set of defined flags.
//
// It wraps flag.FlagSet and provides optional support for the conventional go-service configuration input
// flag ("-i") used by the config subsystem to route decoding.
//
// This type is intentionally minimal: you can still use the embedded *flag.FlagSet to define and parse
// arbitrary flags.
type FlagSet struct {
	input *string
	*flag.FlagSet
}

// AddInput adds the conventional configuration input flag ("-i") to the flag set.
//
// The value is treated as an opaque "kind:location" string that the config package interprets.
// Common examples include:
//   - "file:/path/to/config.yaml" (decode kind inferred from file extension)
//   - "env:MY_CONFIG" (read config payload from environment variable MY_CONFIG)
//
// The provided value is used as the default.
func (f *FlagSet) AddInput(value string) {
	f.input = f.String("i", value, "input config location (format kind:location)")
}

// GetInput returns the configured input flag ("-i") value.
//
// This method is nil-safe with respect to AddInput: if AddInput was not called, GetInput returns an empty
// string. This allows callers to depend on GetInput unconditionally.
func (f *FlagSet) GetInput() string {
	if f.input != nil {
		return *f.input
	}
	return strings.Empty
}
