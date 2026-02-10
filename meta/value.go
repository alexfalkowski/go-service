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

// Ignored constructs a Value that is stored but rendered as an empty string.
//
// This is useful when you want to retain the underlying value in-context, but avoid exporting it
// to logs or transport headers.
func Ignored(value string) Value {
	return Value{kind: ignored, value: value}
}

// Redacted constructs a Value that is rendered as a fixed-length mask of '*' characters.
func Redacted(value string) Value {
	return Value{kind: redacted, value: value}
}

// Blank constructs a Value that represents the absence of a value.
//
// Blank values are treated as empty when exporting attributes to strings.
func Blank() Value {
	return Value{kind: blank, value: strings.Empty}
}

// String constructs a normal Value rendered as-is.
func String(value string) Value {
	return Value{kind: normal, value: value}
}

// Error converts err to a Value using err.Error().
func Error(err error) Value {
	return String(err.Error())
}

// ToString converts st to a normal Value using st.String().
func ToString(st fmt.Stringer) Value {
	return String(st.String())
}

// ToRedacted converts st to a redacted Value using st.String().
func ToRedacted(st fmt.Stringer) Value {
	return Redacted(st.String())
}

// ToIgnored converts st to an ignored Value using st.String().
func ToIgnored(st fmt.Stringer) Value {
	return Ignored(st.String())
}

// Value holds a metadata value along with rendering semantics.
//
// The underlying value can be retrieved with Value.Value(). Rendering is controlled by kind:
//   - normal: renders the value as-is
//   - blank: renders as empty
//   - ignored: renders as empty
//   - redacted: renders as asterisks with the same length as the underlying value
type Value struct {
	value string
	kind  kind
}

// Value returns the underlying value, regardless of rendering semantics.
func (v Value) Value() string {
	return v.value
}

// IsEmpty reports whether the underlying value is empty.
func (v Value) IsEmpty() bool {
	return strings.IsEmpty(v.value)
}

// String renders v according to its kind.
//
//nolint:exhaustive
func (v Value) String() string {
	switch v.kind {
	case redacted:
		return strings.Repeat("*", len(v.value))
	case ignored:
		return strings.Empty
	default:
		return v.value
	}
}
