package config

import (
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/strings"
)

type fileDecoder struct {
	fs       *os.FS
	enc      *encoding.Map
	location string
}

func (f *fileDecoder) Decode(v any) error {
	location := f.fs.CleanPath(f.location)

	file, err := f.fs.Open(location)
	if err != nil {
		return err
	}
	defer file.Close()

	extension := f.fs.PathExtension(location)
	enc := f.enc.Get(extension)
	if enc == nil {
		return errors.Prefix(strings.Join(strings.Space, "file", location, "extension", extension), ErrNoEncoder)
	}

	return enc.Decode(file, v)
}
