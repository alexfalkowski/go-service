package flags

import (
	"github.com/spf13/cobra"
)

// String for cmd.
func String() *string {
	var s string

	return &s
}

// StringVar for cmd.
func StringVar(cmd *cobra.Command, p *string, name, shorthand string, value string, usage string) {
	cmd.PersistentFlags().StringVarP(p, name, shorthand, value, usage)
}

// IsSet the flag for cmd.
func IsSet(b *bool) bool {
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
