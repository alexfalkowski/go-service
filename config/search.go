package config

import (
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/strings"
)

type searchDecoder struct {
	enc  *encoding.Map
	fs   *os.FS
	name env.Name
}

func (c *searchDecoder) Decode(v any) error {
	kind, file, err := c.find()
	if err != nil {
		return err
	}

	defer file.Close()

	return c.enc.Get(kind).Decode(file, v)
}

func (c *searchDecoder) find() (string, io.ReadCloser, error) {
	name := c.name.String()
	dirs := []string{
		c.fs.ExecutableDir(),
		c.fs.Join(os.UserConfigDir(), name),
		"/etc/" + name,
	}
	extensions := []string{".yaml", ".hjson", ".toml", ".json"}
	var paths []string

	for _, extension := range extensions {
		file := name + extension

		for _, dir := range dirs {
			path := c.fs.Join(dir, file)
			paths = append(paths, path)
			if !c.fs.PathExists(path) {
				continue
			}

			f, err := c.fs.Open(path)

			return c.fs.PathExtension(path), f, err
		}
	}

	searched := strings.Join(", ", paths...)
	context := strings.Join(strings.Space, "default config", name, "searched", searched)

	return strings.Empty, nil, errors.Prefix(context, ErrLocationMissing)
}
