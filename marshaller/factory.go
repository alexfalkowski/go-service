package marshaller

import (
	"errors"

	"go.uber.org/fx"
)

// ErrInvalidKind for marshaller.
var ErrInvalidKind = errors.New("invalid kind")

// Factory of marshaller.
type Factory struct {
	yaml *YAML
	toml *TOML
}

// FactoryParams for marshaller.
type FactoryParams struct {
	fx.In

	YAML *YAML
	TOML *TOML
}

// NewFactory for marshaller.
func NewFactory(params FactoryParams) *Factory {
	return &Factory{yaml: params.YAML}
}

// Create from kind.
func (f *Factory) Create(kind string) (Marshaller, error) {
	switch kind {
	case "yaml", "yml":
		return f.yaml, nil
	case "toml":
		return f.toml, nil
	}

	return nil, ErrInvalidKind
}
