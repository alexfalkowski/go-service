package cmd

import (
	"slices"

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

// AddOutput adds an output flag to the flag set.
func (f *FlagSet) AddOutput(value string) {
	f.StringP("output", "o", value, "output config location (format kind:location)")
}

// Provide returns the flag set.
func (f *FlagSet) Provide() *FlagSet {
	return f
}

// SanitizeArgs removes all flags that start with -test.
func SanitizeArgs(args []string) []string {
	return slices.DeleteFunc(args, func(s string) bool {
		return strings.HasPrefix(s, "-test")
	})
}

// SplitFlag for cmd.
func SplitFlag(flag string) (string, string) {
	kind, name, ok := strings.Cut(flag, ":")
	if !ok {
		return "none", "not_used"
	}

	return kind, name
}
