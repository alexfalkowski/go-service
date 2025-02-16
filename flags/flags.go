package flags

import (
	"slices"
	"strings"

	"github.com/spf13/pflag"
)

// FlagSet is a type alias for pflag.FlagSet.
type FlagSet = pflag.FlagSet

// NewFlagSet creates a new flag set with the given name.
func NewFlagSet(name string) *FlagSet {
	return pflag.NewFlagSet(name, pflag.ContinueOnError)
}

// Sanitize removes all flags that start with -test.
func Sanitize(args []string) []string {
	return slices.DeleteFunc(args, func(s string) bool {
		return strings.HasPrefix(s, "-test")
	})
}

// IsStringSet the flag for cmd.
func IsStringSet(s *string) bool {
	return s != nil && *s != ""
}

// String for cmd.
func String() *string {
	var s string

	return &s
}

// IsBoolSet the flag for cmd.
func IsBoolSet(b *bool) bool {
	return b != nil && *b
}

// Bool for cmd.
func Bool() *bool {
	var b bool

	return &b
}
