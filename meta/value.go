package meta

import (
	"fmt"
	"strings"
)

type kind int

const (
	normal kind = iota
	blank
	ignored
	redacted
)

// NewIgnored for meta.
func Ignored(value string) Value {
	return Value{kind: ignored, value: value}
}

// NewRedacted for meta.
func Redacted(value string) Value {
	return Value{kind: redacted, value: value}
}

// NewBlank for meta.
func Blank() Value {
	return Value{kind: blank, value: ""}
}

// NewString for meta.
func String(value string) Value {
	return Value{kind: normal, value: value}
}

// Error for meta.
func Error(err error) Value {
	return String(err.Error())
}

// ToString for meta.
func ToString(st fmt.Stringer) Value {
	return String(st.String())
}

// ToRedacted for meta.
func ToRedacted(st fmt.Stringer) Value {
	return Redacted(st.String())
}

// ToIgnored for meta.
func ToIgnored(st fmt.Stringer) Value {
	return Ignored(st.String())
}

// Value for meta.
type Value struct {
	value string
	kind  kind
}

// Value of meta.
func (v Value) Value() string {
	return v.value
}

// String depending on kind.
//
//nolint:exhaustive
func (v Value) String() string {
	switch v.kind {
	case redacted:
		return strings.Repeat("*", len(v.value))
	case ignored:
		return ""
	default:
		return v.value
	}
}

// IsBlank checks for an empty string.
func (v Value) IsBlank() bool {
	return v.IsEqual("")
}

// IsEqual is true if values match.
func (v Value) IsEqual(value string) bool {
	return v.value == value
}
