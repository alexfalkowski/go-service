package bytes

import (
	"github.com/alexfalkowski/go-service/v2/encoding/json"
	"github.com/alexfalkowski/go-service/v2/runtime"
	units "github.com/docker/go-units"
)

// KB is the decimal kilobyte size constant: 1,000 bytes.
const KB Size = Size(units.KB)

// MB is the decimal megabyte size constant: 1,000,000 bytes.
const MB Size = Size(units.MB)

// GB is the decimal gigabyte size constant: 1,000,000,000 bytes.
const GB Size = Size(units.GB)

// TB is the decimal terabyte size constant: 1,000,000,000,000 bytes.
const TB Size = Size(units.TB)

// PB is the decimal petabyte size constant: 1,000,000,000,000,000 bytes.
const PB Size = Size(units.PB)

// Size is the go-service decimal byte-size type used across the repository.
//
// It is a named type over int64 so it can expose config-friendly text and JSON
// marshaling helpers while remaining easy to convert at API boundaries.
//
// Size uses the SI units understood by `github.com/docker/go-units`, such as
// `B`, `kB`, `MB`, `GB`, and `TB`, rather than IEC binary units like `MiB`.
type Size int64

// Bytes returns s as a raw byte count.
func (s Size) Bytes() int64 {
	return int64(s)
}

// String returns s in the human-readable decimal size format used by this
// package, such as `64B` or `4MB`.
func (s Size) String() string {
	return units.HumanSize(float64(s.Bytes()))
}

// MarshalText encodes s using the same decimal size string returned by
// [Size.String].
func (s Size) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

// UnmarshalText parses a decimal size string into s.
//
// Accepted inputs use the same format as [ParseSize], such as `64B`, `2MB`, or
// `4GB`.
func (s *Size) UnmarshalText(text []byte) error {
	size, err := ParseSize(string(text))
	if err != nil {
		return err
	}

	*s = size
	return nil
}

// MarshalJSON encodes s as a quoted decimal size string, such as `"4MB"`.
func (s Size) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// UnmarshalJSON decodes a quoted decimal size string into s.
//
// Non-string JSON values are rejected by the underlying JSON decoder before the
// size parser runs.
func (s *Size) UnmarshalJSON(data []byte) error {
	var text string
	if err := json.Unmarshal(data, &text); err != nil {
		return err
	}

	return s.UnmarshalText([]byte(text))
}

// ParseSize parses a human-readable SI size string.
//
// The input uses decimal SI units such as `64B`, `2MB`, or `4GB`. Parsing is
// delegated to `github.com/docker/go-units`.
func ParseSize(s string) (Size, error) {
	size, err := units.FromHumanSize(s)
	return Size(size), err
}

// MustParseSize parses s as a human-readable SI size string and panics if parsing fails.
//
// This helper is intended for strict startup/configuration paths where an invalid
// size is considered a fatal configuration/programming error. It panics by
// calling runtime.Must on the parse error.
//
// If you need recoverable error handling, use ParseSize instead.
func MustParseSize(s string) Size {
	size, err := ParseSize(s)
	runtime.Must(err)
	return size
}
