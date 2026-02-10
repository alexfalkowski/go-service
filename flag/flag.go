package flag

import (
	"flag"

	"github.com/alexfalkowski/go-service/v2/strings"
)

// NewFlagSet creates a new FlagSet with the given name.
func NewFlagSet(name string) *FlagSet {
	set := flag.NewFlagSet(name, flag.ContinueOnError)
	return &FlagSet{FlagSet: set}
}

// FlagSet represents a set of defined flags.
//
// It wraps flag.FlagSet and may optionally include an input configuration flag ("-i") for selecting a config source.
type FlagSet struct {
	input *string
	*flag.FlagSet
}

// AddInput adds an input config source flag ("-i") to the flag set.
//
// The value format is "kind:location", for example:
//   - "file:/path/to/config.yaml"
//   - "env:MY_CONFIG"
func (f *FlagSet) AddInput(value string) {
	f.input = f.String("i", value, "input config location (format kind:location)")
}

// GetInput retrieves the input flag value from the flag set.
//
// If AddInput was not called, GetInput returns an empty string.
func (f *FlagSet) GetInput() string {
	if f.input != nil {
		return *f.input
	}
	return strings.Empty
}
