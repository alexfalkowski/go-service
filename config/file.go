package config

import (
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/os"
)

// NewFile constructs a Decoder that loads configuration from a specific file path.
//
// The decoder selects an encoder/decoder based on the file extension of location.
func NewFile(location string, enc *encoding.Map, fs *os.FS) *File {
	return &File{location: location, enc: enc, fs: fs}
}

// File decodes configuration from a specific file path.
type File struct {
	fs       *os.FS
	enc      *encoding.Map
	location string
}

// Decode opens the configured file and decodes its contents into v.
//
// It returns ErrNoEncoder when no encoder is registered for the file extension.
// It returns any filesystem open errors and any decode errors returned by the selected encoder.
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
