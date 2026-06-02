package bytes

import (
	"strconv"

	"github.com/alexfalkowski/go-service/v2/runtime"
	units "github.com/docker/go-units"
)

const (
	// KB is the decimal kilobyte size constant: 1,000 bytes.
	KB Size = Size(units.KB)

	// MB is the decimal megabyte size constant: 1,000,000 bytes.
	MB Size = Size(units.MB)

	// GB is the decimal gigabyte size constant: 1,000,000,000 bytes.
	GB Size = Size(units.GB)

	// TB is the decimal terabyte size constant: 1,000,000,000,000 bytes.
	TB Size = Size(units.TB)

	// PB is the decimal petabyte size constant: 1,000,000,000,000,000 bytes.
	PB Size = Size(units.PB)
)

const (
	// DefaultSize is the shared default size used by config surfaces that need a conservative byte limit.
	//
	// Its value is 4 megabytes.
	DefaultSize Size = 4 * MB

	// MaxConfigSize is the largest byte size accepted by configuration surfaces.
	//
	// Configured sizes protect in-memory buffering paths such as cache values, HTTP request bodies,
	// HTTP client response bodies, and gRPC messages. Keep this comfortably below integer and allocator
	// edge cases while still allowing unusually large service payloads.
	MaxConfigSize Size = 256 * MB
)

// Size is the go-service decimal byte-size type used across the repository.
//
// It is a named type over int64 so it can expose config-friendly text and JSON
// marshaling helpers while remaining easy to convert at API boundaries.
//
// Size uses the decimal units understood by `github.com/docker/go-units`, such
// as `B`, `kB`, `MB`, `GB`, and `TB`.
//
// Formatting and parsing intentionally follow go-units compatibility behavior.
// They are not a strict inverse for exabyte-scale values because go-units can
// format larger suffixes than FromHumanSize parses.
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

// MarshalText encodes s as an exact raw byte count with a `B` suffix.
func (s Size) MarshalText() ([]byte, error) {
	return []byte(strconv.FormatInt(s.Bytes(), 10) + "B"), nil
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

// MarshalJSON encodes s as a quoted exact raw byte count, such as `"4000000B"`.
func (s Size) MarshalJSON() ([]byte, error) {
	text, err := s.MarshalText()
	if err != nil {
		return nil, err
	}

	return []byte(strconv.Quote(string(text))), nil
}

// UnmarshalJSON decodes a quoted decimal size string into s.
//
// Non-string JSON values are rejected before the size parser runs.
func (s *Size) UnmarshalJSON(data []byte) error {
	data = TrimSpace(data)
	if len(data) == 0 || data[0] != '"' {
		return strconv.ErrSyntax
	}

	text, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}

	return s.UnmarshalText([]byte(text))
}

// ParseSize parses a human-readable decimal size string.
//
// The input uses the suffix compatibility rules from
// `github.com/docker/go-units.FromHumanSize`. Parsed values are decimal sizes,
// so accepted suffix spellings such as `MB` and `MiB` both use the decimal `M`
// multiplier.
func ParseSize(s string) (Size, error) {
	size, err := units.FromHumanSize(s)
	return Size(size), err
}

// MustParseSize parses s as a human-readable decimal size string and panics if parsing fails.
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
