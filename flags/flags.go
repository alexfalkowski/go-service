package flags

import (
	"slices"
	"strings"

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

// AddOutput adds an output flag to the flag set.
func (f *FlagSet) AddOutput(value string) {
	f.StringP("output", "o", value, "output config location (format kind:location)")
}

// Flags returns the flag set.
func (f *FlagSet) Flags() *FlagSet {
	return f
}

// Sanitize removes all flags that start with -test.
func Sanitize(args []string) []string {
	return slices.DeleteFunc(args, func(s string) bool {
		return strings.HasPrefix(s, "-test")
	})
}
