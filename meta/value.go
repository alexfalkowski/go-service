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

// Ignored constructs a Value that preserves the underlying value in-context but renders as an empty string.
//
// This is useful when you want to keep the raw value available for in-process logic, but you do not want it
// to be exported via `meta.Strings` / `meta.SnakeStrings` / `meta.CamelStrings` (for example into logs or
// transport headers).
//
// Note: export helpers skip attributes whose rendered string is empty, so Ignored values will not appear
// in exported maps.
func Ignored(value string) Value {
	return Value{kind: ignored, value: value}
}

// Redacted constructs a Value that renders as a fixed-length mask of '*' characters.
//
// The rendered value has the same length as the underlying value, which can be useful to preserve the
// shape of data without revealing the contents.
//
// Note: the underlying value is still stored in-context and can be retrieved via Value.Value(). Use this
// only when it is acceptable for in-process code to access the raw value.
func Redacted(value string) Value {
	return Value{kind: redacted, value: value}
}

// Blank constructs a Value that represents the absence of a value.
//
// Blank values render as empty strings and are skipped by export helpers.
//
// Use Blank when you want to explicitly represent "no value" rather than "an ignored value that still exists".
func Blank() Value {
	return Value{kind: blank, value: strings.Empty}
}

// String constructs a normal Value that renders as-is.
func String(value string) Value {
	return Value{kind: normal, value: value}
}

// Error converts err to a normal Value using err.Error().
//
// Callers should ensure err is non-nil before calling Error; passing nil will panic.
func Error(err error) Value {
	return String(err.Error())
}

// ToString converts st to a normal Value using st.String().
//
// Callers should ensure st is non-nil before calling ToString; passing nil will panic.
func ToString(st fmt.Stringer) Value {
	return String(st.String())
}

// ToRedacted converts st to a redacted Value using st.String().
//
// Callers should ensure st is non-nil before calling ToRedacted; passing nil will panic.
// The underlying (unredacted) string is still stored in-context and can be retrieved via Value.Value().
func ToRedacted(st fmt.Stringer) Value {
	return Redacted(st.String())
}

// ToIgnored converts st to an ignored Value using st.String().
//
// Callers should ensure st is non-nil before calling ToIgnored; passing nil will panic.
// The underlying string is still stored in-context and can be retrieved via Value.Value().
func ToIgnored(st fmt.Stringer) Value {
	return Ignored(st.String())
}

// Value holds a metadata value along with rendering semantics.
//
// Value is designed to support two related use-cases:
//
//  1. In-process access: retrieve the underlying value with Value.Value() for business logic.
//  2. Export: render the value with Value.String() when exporting metadata (logs/headers/etc.).
//
// Rendering is controlled by kind:
//   - normal: renders the underlying value as-is
//   - blank: renders as empty
//   - ignored: renders as empty (while still retaining the underlying value in-context)
//   - redacted: renders as asterisks with the same length as the underlying value
//
// Note: Value.String intentionally does not distinguish between blank and ignored during rendering; both
// render as empty. The difference is semantic at construction time.
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
// For redacted values, the returned string is a mask of '*' characters with the same length as the underlying
// value. For ignored values, it returns an empty string. For normal and blank values, it returns the underlying
// value (blank values use an empty underlying string).
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
