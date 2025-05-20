package flag

import (
	"github.com/alexfalkowski/go-service/strings"
	"github.com/spf13/pflag"
)

// NewFlagSet creates a new flag set with the given name.
func NewFlagSet(name string) *FlagSet {
	set := pflag.NewFlagSet(name, pflag.ContinueOnError)

	return &FlagSet{FlagSet: set}
}

// FlagSet represents a set of defined flags.
type FlagSet struct {
	*pflag.FlagSet
}

// AddInput adds an input flag to the flag set.
func (f *FlagSet) AddInput(value string) {
	f.StringP("input", "i", value, "input config location (format kind:location)")
}

// GetInput retrieves the input flag value from the flag set.
func (f *FlagSet) GetInput() string {
	input, _ := f.GetString("input")

	return input
}

// AddOutput adds an output flag to the flag set.
func (f *FlagSet) AddOutput(value string) {
	f.StringP("output", "o", value, "output config location (format kind:location)")
}

// GetInput retrieves the output flag value from the flag set.
func (f *FlagSet) GetOutput() string {
	output, _ := f.GetString("output")

	return output
}

// Provide returns the flag set.
func (f *FlagSet) Provide() *FlagSet {
	return f
}

// SplitFlag will split by : and give us the kind and name.
func SplitFlag(arg string) (string, string) {
	kind, name, ok := strings.Cut(arg, ":")
	if !ok {
		return "", ""
	}

	return kind, name
}
