package meta

import (
	"fmt"

	"github.com/alexfalkowski/go-service/v2/strings"
)

type kind uint8

const (
	normal kind = iota
	blank
	ignored
	redacted
)

// Ignored for meta.
func Ignored(value string) Value {
	return Value{kind: ignored, value: value}
}

// Redacted for meta.
func Redacted(value string) Value {
	return Value{kind: redacted, value: value}
}

// Blank for meta.
func Blank() Value {
	return Value{kind: blank, value: ""}
}

// String for meta.
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

// IsEmpty of meta.
func (v Value) IsEmpty() bool {
	return strings.IsEmpty(v.value)
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
