package config

import (
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/os"
)

// NewCommon for config.
func NewCommon(name env.Name, enc *encoding.Map, fs *os.FS) *Common {
	return &Common{name: name, enc: enc, fs: fs}
}

// Common for config.
type Common struct {
	enc  *encoding.Map
	fs   *os.FS
	name env.Name
}

// Decode to v.
func (c *Common) Decode(v any) error {
	kind, file, err := c.find()
	if err != nil {
		return err
	}

	defer file.Close()

	return c.enc.Get(kind).Decode(file, v)
}

func (c *Common) find() (string, io.ReadCloser, error) {
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

	return "", nil, ErrLocationMissing
}
