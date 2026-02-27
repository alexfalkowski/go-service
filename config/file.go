package config

import (
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/os"
)

// NewFile constructs a file-based Decoder that loads configuration from a specific file path.
//
// The decoder selects an encoder based on the file extension of location (for example ".yaml" → "yaml").
// If no encoder is registered for the extension, Decode will return ErrNoEncoder.
//
// Note: location is used as-is when opening the file; any path expansion/cleaning is the responsibility
// of the underlying filesystem implementation.
func NewFile(location string, enc *encoding.Map, fs *os.FS) *File {
	return &File{location: location, enc: enc, fs: fs}
}

// File decodes configuration from a specific file path.
//
// The file extension determines which encoder is used (via encoding.Map). This decoder does not attempt
// to infer the configuration kind from content.
type File struct {
	fs       *os.FS
	enc      *encoding.Map
	location string
}

// Decode opens the configured file and decodes its contents into v.
//
// The destination v should be a pointer to the target configuration type.
//
// Errors:
//   - returns filesystem errors if the file cannot be opened,
//   - returns ErrNoEncoder if no encoder is registered for the file extension, and
//   - returns any decode/unmarshal error from the selected encoder.
func (f *File) Decode(v any) error {
	file, err := f.fs.Open(f.location)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := f.enc.Get(f.fs.PathExtension(f.location))
	if enc == nil {
		return ErrNoEncoder
	}

	return enc.Decode(file, v)
}
