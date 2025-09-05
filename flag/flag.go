package flag

import "flag"

// NewFlagSet creates a new flag set with the given name.
func NewFlagSet(name string) *FlagSet {
	set := flag.NewFlagSet(name, flag.ContinueOnError)
	return &FlagSet{FlagSet: set}
}

// FlagSet represents a set of defined flags.
type FlagSet struct {
	input *string
	*flag.FlagSet
}

// AddInput adds an input flag to the flag set.
func (f *FlagSet) AddInput(value string) {
	f.input = f.String("i", value, "input config location (format kind:location)")
}

// GetInput retrieves the input flag value from the flag set.
func (f *FlagSet) GetInput() string {
	if f.input != nil {
		return *f.input
	}
	return ""
}
