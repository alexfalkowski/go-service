package meta

import (
	"fmt"
	"strings"
)

// Valuer for meta.
type Valuer interface {
	Value() string

	fmt.Stringer
}

// ValueOrBlank for meta.
func ValueOrBlank(s Valuer) string {
	if s == nil {
		return ""
	}

	return s.Value()
}

// IsBlank checks for an empty string.
func IsBlank(actual Valuer) bool {
	return IsEqual(actual, "")
}

// IsEqual checks if actual is expected.
func IsEqual(actual Valuer, expected string) bool {
	return ValueOrBlank(actual) == expected
}

// Error for meta.
func Error(err error) Valuer {
	if err == nil {
		return String("")
	}

	return String(err.Error())
}

// ToString for meta.
func ToString(st fmt.Stringer) String {
	if st == nil {
		return String("")
	}

	return String(st.String())
}

// String for meta.
type String string

// Value of the string.
func (v String) Value() string {
	return string(v)
}

// String to satisfy fmt.Stringer.
func (v String) String() string {
	return v.Value()
}

// ToRedacted for meta.
func ToRedacted(st fmt.Stringer) Redacted {
	if st == nil {
		return Redacted("")
	}

	return Redacted(st.String())
}

// Redacted for meta.
type Redacted string

// Value of the string.
func (v Redacted) Value() string {
	return string(v)
}

// String to satisfy fmt.Stringer.
func (v Redacted) String() string {
	return strings.Repeat("*", len(v.Value()))
}
