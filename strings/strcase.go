package strings

import "github.com/iancoleman/strcase"

// ToDelimited converts s to a delimiter-separated string.
//
// This is a thin wrapper around [strcase.ToDelimited] and does not change semantics.
func ToDelimited(s string, delimiter uint8) string {
	return strcase.ToDelimited(s, delimiter)
}

// ToLowerCamel converts s to lower camel case.
//
// This is a thin wrapper around [strcase.ToLowerCamel] and does not change semantics.
func ToLowerCamel(s string) string {
	return strcase.ToLowerCamel(s)
}

// ToSnake converts s to snake case.
//
// This is a thin wrapper around [strcase.ToSnake] and does not change semantics.
func ToSnake(s string) string {
	return strcase.ToSnake(s)
}
