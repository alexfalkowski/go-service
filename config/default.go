package config

import (
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// NewDefault constructs a Decoder that locates a configuration file by searching common locations.
//
// The returned decoder uses the provided service name to look for "<serviceName>.<ext>" where <ext> is one
// of the supported extensions (see (*Default).Decode).
func NewDefault(name env.Name, enc *encoding.Map, fs *os.FS) *Default {
	return &Default{name: name, enc: enc, fs: fs}
}

// Default is a Decoder implementation that performs "default lookup" configuration discovery.
//
// It searches for a config file named "<serviceName>.<ext>" in a fixed set of directories and decodes it
// using the encoder registered for that extension.
type Default struct {
	enc  *encoding.Map
	fs   *os.FS
	name env.Name
}

// Decode decodes configuration into v using default lookup.
//
// It searches for the first existing file named "<serviceName>.<ext>", where <ext> is one of:
//   - .yaml
//   - .yml
//   - .toml
//   - .json
//
// For each extension, directories are checked in order:
//   - the executable directory (fs.ExecutableDir())
//   - the user config directory under "<configDir>/<serviceName>/" (os.UserConfigDir())
//   - /etc/<serviceName>/
//
// The first match is opened and decoded using the encoder keyed by the discovered kind/extension.
// If no configuration file is found, Decode returns ErrLocationMissing.
//
// Note: this decoder assumes the encoding.Map contains decoders for the supported kinds. If a kind is
// missing, Decode may panic due to a nil encoder; ensure your wiring registers the standard encoders.
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
