package compressor

import (
	"errors"
)

type configs map[string]Compressor

// ErrInvalidKind for compressor.
var ErrInvalidKind = errors.New("invalid kind")

// Factory of compressor.
type Factory struct {
	configs configs
}

// NewFactory for compressor.
func NewFactory() *Factory {
	f := &Factory{
		configs: configs{
			"snappy": NewSnappy(),
			"none":   NewNone(),
		},
	}

	return f
}

// Register kind and compressor.
func (f *Factory) Register(kind string, c Compressor) {
	f.configs[kind] = c
}

// Create from kind.
func (f *Factory) Create(kind string) (Compressor, error) {
	c, ok := f.configs[kind]
	if !ok {
		return nil, ErrInvalidKind
	}

	return c, nil
}
