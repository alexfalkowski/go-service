package bytes

import (
	"encoding/json"

	"github.com/alexfalkowski/go-service/v2/runtime"
	units "github.com/docker/go-units"
)

// KB is a size constant equal to 1000 bytes.
const KB Size = Size(units.KB)

// MB is a size constant equal to 1000 kilobytes.
const MB Size = Size(units.MB)

// GB is a size constant equal to 1000 megabytes.
const GB Size = Size(units.GB)

// TB is a size constant equal to 1000 gigabytes.
const TB Size = Size(units.TB)

// PB is a size constant equal to 1000 terabytes.
const PB Size = Size(units.PB)

// Size is the go-service byte-size type used across the repository.
//
// It is a named type over int64 so it can expose config-friendly marshaling
// helpers while remaining easy to convert at API boundaries.
type Size int64

// Bytes converts s to its raw byte count.
func (s Size) Bytes() int64 {
	return int64(s)
}

// String returns the human-readable SI size string for s.
func (s Size) String() string {
	return units.HumanSize(float64(s.Bytes()))
}

// MarshalText encodes s using the human-readable SI size format.
func (s Size) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

// UnmarshalText decodes a human-readable SI size string into s.
func (s *Size) UnmarshalText(text []byte) error {
	size, err := ParseSize(string(text))
	if err != nil {
		return err
	}

	*s = size
	return nil
}

// MarshalJSON encodes s as a quoted human-readable SI size string.
func (s Size) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// UnmarshalJSON decodes a quoted human-readable SI size string into s.
func (s *Size) UnmarshalJSON(data []byte) error {
	var text string
	if err := json.Unmarshal(data, &text); err != nil {
		return err
	}

	return s.UnmarshalText([]byte(text))
}

// ParseSize parses a human-readable SI size string.
//
// The input uses SI units such as "64B", "2MB", or "4GB".
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
