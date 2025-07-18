package config

import (
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/os"
)

// NewFile for config.
func NewFile(location string, enc *encoding.Map, fs *os.FS) *File {
	return &File{location: location, enc: enc, fs: fs}
}

// File for config.
type File struct {
	fs       *os.FS
	enc      *encoding.Map
	location string
}

// Decode to v.
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
