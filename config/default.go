package config

import (
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// NewDefault for config.
func NewDefault(name env.Name, enc *encoding.Map, fs *os.FS) *Default {
	return &Default{name: name, enc: enc, fs: fs}
}

// Default for config.
type Default struct {
	enc  *encoding.Map
	fs   *os.FS
	name env.Name
}

// Decode to v, by looking in default places.
func (c *Default) Decode(v any) error {
	kind, file, err := c.find()
	if err != nil {
		return err
	}

	defer file.Close()

	return c.enc.Get(kind).Decode(file, v)
}

func (c *Default) find() (string, io.ReadCloser, error) {
	extensions := []string{".yaml", ".yml", ".toml", ".json"}
	for _, extension := range extensions {
		n := c.name.String()
		file := n + extension
		dirs := []string{
			c.fs.ExecutableDir(),
			c.fs.Join(os.UserConfigDir(), n),
			"/etc/" + n,
		}

		for _, dir := range dirs {
			name := c.fs.Join(dir, file)
			if !c.fs.PathExists(name) {
				continue
			}

			f, err := c.fs.Open(name)

			return c.fs.PathExtension(name), f, err
		}
	}

	return strings.Empty, nil, ErrLocationMissing
}
