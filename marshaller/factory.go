package marshaller

import (
	"errors"
)

type configs map[string]Marshaller

// ErrInvalidKind for marshaller.
var ErrInvalidKind = errors.New("invalid kind")

// Factory of marshaller.
type Factory struct {
	configs configs
}

// NewFactory for marshaller.
func NewFactory() *Factory {
	f := &Factory{
		configs: configs{
			"json":  NewJSON(),
			"yaml":  NewYAML(),
			"yml":   NewYAML(),
			"toml":  NewTOML(),
			"proto": NewProto(),
		},
	}

	return f
}

// Register kind and marshaller.
func (f *Factory) Register(kind string, m Marshaller) {
	f.configs[kind] = m
}

// Create from kind.
func (f *Factory) Create(kind string) (Marshaller, error) {
	c, ok := f.configs[kind]
	if !ok {
		return nil, ErrInvalidKind
	}

	return c, nil
}
