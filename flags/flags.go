package flags

import (
	"github.com/spf13/cobra"
)

// IsStringSet the flag for cmd.
func IsStringSet(s *string) bool {
	return s != nil && *s != ""
}

// String for cmd.
func String() *string {
	var s string

	return &s
}

// StringVar for cmd.
func StringVar(cmd *cobra.Command, p *string, name, shorthand string, value string, usage string) {
	cmd.PersistentFlags().StringVarP(p, name, shorthand, value, usage)
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

// BoolVar for cmd.
func BoolVar(cmd *cobra.Command, p *bool, name, shorthand string, value bool, usage string) {
	cmd.PersistentFlags().BoolVarP(p, name, shorthand, value, usage)
}
